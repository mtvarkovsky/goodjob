package processor

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/pkg/goodjob"
)

type (
	V1 struct {
		Queue             goodjob.Queue
		JobResultStorage  goodjob.JobResultsStorage
		TaskResultStorage goodjob.TaskResultsStorage
		Active            bool
	}
)

var _ goodjob.Processor = (*V1)(nil)

func NewV1Processor(
	queue goodjob.Queue,
	jobResultsStorage goodjob.JobResultsStorage,
	taskResultsStorage goodjob.TaskResultsStorage,
) goodjob.Processor {
	return &V1{
		Queue:             queue,
		JobResultStorage:  jobResultsStorage,
		TaskResultStorage: taskResultsStorage,
		Active:            false,
	}
}

func (p *V1) Start(args ...*goodjob.ProcessorArg) error {
	p.Active = true
	go func() {
		for p.Active {
			go func() {
				err := p.processNextJob()
				if err != nil {
					// TODO: insert error reporting system
					return
				}
			}()
		}
	}()
	return nil
}

func (p *V1) Stop(args ...*goodjob.ProcessorArg) error {
	p.Active = false
	return nil
}

func (p *V1) AddJob(job goodjob.Job, args ...*goodjob.ProcessorArg) error {
	queueArgs := make([]*goodjob.QueueArg, len(args))
	for i, a := range args {
		queueArgs[i].Name = a.Name
		queueArgs[i].Value = a.Value
	}
	return p.Queue.AddJob(job, queueArgs...)
}

func (p *V1) GetJobResult(id goodjob.JobID, args ...*goodjob.ProcessorArg) (*goodjob.JobResult, error) {
	return p.JobResultStorage.Get(id)
}

func (p *V1) processNextJob() error {
	job, err := p.Queue.GetNextJob()
	if err != nil {
		return err
	}
	if job == nil {
		return nil
	}

	return p.runJob(job)
}

func (p *V1) runJob(job goodjob.Job) error {
	switch j := job.(type) {
	case goodjob.RetryableRevertibleJob:
		p.runRetryableRevertibleJob(j)
		return nil
	case goodjob.RevertibleJob:
		p.runRevertibleJob(j)
		return nil
	case goodjob.RetryableJob:
		p.runRetryableJob(j)
		return nil
	case goodjob.Job:
		p.runBasicJob(j)
		return nil
	}

	return fmt.Errorf("unknown job type")
}

func (p *V1) runBasicJob(job goodjob.Job) {
	tasks := job.GetTasks()

	for i, task := range tasks {
		taskArgs, err := p.prepareTaskArgs(job, task)
		if err != nil {
			// TODO: insert error reporting system
			return
		}

		res, err := p.execTask(job, i, task, taskArgs...)
		if err != nil {
			// TODO: insert error reporting system
			return
		}

		if res.Err != nil {
			break
		}
	}

	err := p.saveJobResult(job)
	if err != nil {
		// TODO: insert error reporting system
		return
	}
	err = p.Queue.RemoveJob(job.GetID())
	if err != nil {
		// TODO: insert error reporting system
		return
	}
}

func (p *V1) prepareTaskArgs(job goodjob.Job, task goodjob.Task) ([]*goodjob.TaskArg, error) {
	taskArgs := job.GetTaskArgs(task.GetID())
	for _, arg := range taskArgs {
		if arg.ValueFrom != "" {
			prevRes, err := p.TaskResultStorage.Get(job.GetID(), arg.ValueFrom)
			if err != nil {
				return nil, err
			}
			arg.Value = prevRes.Value
		}
	}
	return taskArgs, nil
}

func (p *V1) execTask(job goodjob.Job, taskPos int, task goodjob.Task, args ...*goodjob.TaskArg) (*goodjob.TaskResult, error) {
	res := task.Exec(args...)
	err := p.saveTaskResult(job, task, taskPos, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *V1) revertTask(job goodjob.Job, taskPos int, task goodjob.RevertibleTask, args ...*goodjob.TaskArg) (*goodjob.TaskResult, error) {
	res := task.Revert(args...)
	err := p.saveTaskResult(job, task, taskPos, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *V1) saveTaskResult(job goodjob.Job, task goodjob.Task, taskPos int, res *goodjob.TaskResult) error {
	err := p.TaskResultStorage.Put(res)
	if err != nil {
		return err
	}
	job.SetLastTask(task)
	job.SetLastTaskPos(taskPos)
	job.SetLastTaskResult(res)

	return nil
}

func (p *V1) saveJobResult(job goodjob.Job) error {
	lastResult := job.GetLastTaskResult()
	err := p.JobResultStorage.Put(&goodjob.JobResult{
		ID:    lastResult.JobID,
		Value: lastResult.Value,
		Err:   lastResult.Err,
	})
	return err
}

func (p *V1) runRetryableJob(job goodjob.RetryableJob) {
	var tasks []goodjob.Task

	lastRunTask := job.GetLastTask()

	if job.GetLastTaskResult() != nil && job.GetLastTaskResult().Err != nil {
		job.IncreaseRetryCount(lastRunTask.GetID())
		tasks = job.GetTasksToRetry()
	} else {
		tasks = job.GetTasks()
	}

	lastPos := job.GetLastTaskPos()

	for i, task := range tasks {
		taskArgs, err := p.prepareTaskArgs(job, task)
		if err != nil {
			// TODO: insert error reporting system
			return
		}

		res, err := p.execTask(job, lastPos+i, task, taskArgs...)
		if err != nil {
			// TODO: insert error reporting system
			return
		}

		if res.Err != nil {
			if !job.RetryThresholdReached(task.GetID()) {
				err := p.Queue.SetJobVisibility(job.GetID(), true)
				if err != nil {
					// TODO: insert error reporting system
					return
				}
				return
			}
			break
		}
	}

	err := p.saveJobResult(job)
	if err != nil {
		// TODO: insert error reporting system
		return
	}
	err = p.Queue.RemoveJob(job.GetID())
	if err != nil {
		// TODO: insert error reporting system
		return
	}
}

func (p *V1) runRevertibleJob(job goodjob.RevertibleJob) {
	var tasks []goodjob.Task
	var tasksToRevert []goodjob.RevertibleTask

	if job.GetRevertState() {
		tasksToRevert = job.GetTasksToRevert()
	} else {
		tasks = job.GetTasks()
	}

	for i, task := range tasks {
		taskArgs, err := p.prepareTaskArgs(job, task)
		if err != nil {
			// TODO: insert error reporting system
			return
		}

		res, err := p.execTask(job, i, task, taskArgs...)
		if err != nil {
			// TODO: insert error reporting system
			return
		}

		if res.Err != nil {
			job.SetRevertState(true)
			err := p.Queue.SetJobVisibility(job.GetID(), true)
			if err != nil {
				// TODO: insert error reporting system
				return
			}
			return
		}
	}

	lastPos := job.GetLastTaskPos()

	for i, task := range tasksToRevert {
		taskArgs, err := p.prepareTaskArgs(job, task)
		if err != nil {
			// TODO: insert error reporting system
			return
		}

		res, err := p.revertTask(job, lastPos-i, task, taskArgs...)
		if err != nil {
			// TODO: insert error reporting system
			return
		}

		if res.Err != nil {
			break
		}
	}

	err := p.saveJobResult(job)
	if err != nil {
		// TODO: insert error reporting system
		return
	}
	err = p.Queue.RemoveJob(job.GetID())
	if err != nil {
		// TODO: insert error reporting system
		return
	}
}

func (p *V1) runRetryableRevertibleJob(job goodjob.RetryableRevertibleJob) {
	var tasks []goodjob.Task
	var tasksToRevert []goodjob.RevertibleTask

	lastRunTask := job.GetLastTask()

	if job.GetRevertState() {
		tasksToRevert = job.GetTasksToRevert()
	} else if job.GetLastTaskResult() != nil && job.GetLastTaskResult().Err != nil {
		tasks = job.GetTasksToRetry()
	} else {
		tasks = job.GetTasks()
	}

	if job.GetLastTaskResult() != nil && job.GetLastTaskResult().Err != nil {
		job.IncreaseRetryCount(lastRunTask.GetID())
	}

	lastPos := job.GetLastTaskPos()

	for i, task := range tasks {
		taskArgs, err := p.prepareTaskArgs(job, task)
		if err != nil {
			// TODO: insert error reporting system
			return
		}

		res, err := p.execTask(job, lastPos+i, task, taskArgs...)
		if err != nil {
			// TODO: insert error reporting system
			return
		}

		if res.Err != nil {
			if !job.RetryThresholdReached(task.GetID()) {
				err := p.Queue.SetJobVisibility(job.GetID(), true)
				if err != nil {
					// TODO: insert error reporting system
					return
				}
				return
			} else {
				job.SetRevertState(true)
				err := p.Queue.SetJobVisibility(job.GetID(), true)
				if err != nil {
					// TODO: insert error reporting system
					return
				}
				return
			}
		}
	}

	for i, task := range tasksToRevert {
		taskArgs, err := p.prepareTaskArgs(job, task)
		if err != nil {
			// TODO: insert error reporting system
			return
		}

		res, err := p.revertTask(job, lastPos-i, task, taskArgs...)
		if err != nil {
			// TODO: insert error reporting system
			return
		}

		if res.Err != nil {
			if !job.RetryThresholdReached(task.GetID()) {
				err := p.Queue.SetJobVisibility(job.GetID(), true)
				if err != nil {
					// TODO: insert error reporting system
					return
				}
				return
			}
			break
		}
	}

	err := p.saveJobResult(job)
	if err != nil {
		// TODO: insert error reporting system
		return
	}
	err = p.Queue.RemoveJob(job.GetID())
	if err != nil {
		// TODO: insert error reporting system
		return
	}
}
