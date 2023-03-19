package processor

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/pkg/interfaces"
)

type (
	V1 struct {
		Queue             interfaces.Queue
		JobResultStorage  interfaces.JobResultsStorage
		TaskResultStorage interfaces.TaskResultsStorage
		Active            bool
	}
)

var _ interfaces.Processor = (*V1)(nil)

func NewV1Processor(
	queue interfaces.Queue,
	jobResultsStorage interfaces.JobResultsStorage,
	taskResultsStorage interfaces.TaskResultsStorage,
) interfaces.Processor {
	return &V1{
		Queue:             queue,
		JobResultStorage:  jobResultsStorage,
		TaskResultStorage: taskResultsStorage,
		Active:            false,
	}
}

func (p *V1) Start(args ...any) error {
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

func (p *V1) Stop(args ...any) error {
	p.Active = false
	return nil
}

func (p *V1) AddJob(job interfaces.Job, args ...any) error {
	return p.Queue.AddJob(job, args...)
}

func (p *V1) GetJobResult(id interfaces.JobID, args ...any) (*interfaces.JobResult, error) {
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

func (p *V1) runJob(job interfaces.Job) error {
	switch j := job.(type) {
	case interfaces.RetryableRevertibleJob:
		p.runRetryableJob(j)
		return nil
	case interfaces.RevertibleJob:
		p.runRevertibleJob(j)
		return nil
	case interfaces.RetryableJob:
		p.runRetryableJob(j)
		return nil
	case interfaces.Job:
		p.runBasicJob(j)
		return nil
	}

	return fmt.Errorf("unknown job type")
}

func (p *V1) runBasicJob(job interfaces.Job) {
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
}

func (p *V1) prepareTaskArgs(job interfaces.Job, task interfaces.Task) ([]*interfaces.TaskArg, error) {
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

func (p *V1) execTask(job interfaces.Job, taskPos int, task interfaces.Task, args ...*interfaces.TaskArg) (*interfaces.TaskResult, error) {
	res := task.Exec(args...)
	err := p.saveTaskResult(job, task, taskPos, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *V1) revertTask(job interfaces.Job, taskPos int, task interfaces.RevertibleTask, args ...*interfaces.TaskArg) (*interfaces.TaskResult, error) {
	res := task.Revert(args...)
	err := p.saveTaskResult(job, task, taskPos, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *V1) saveTaskResult(job interfaces.Job, task interfaces.Task, taskPos int, res *interfaces.TaskResult) error {
	err := p.TaskResultStorage.Put(res)
	if err != nil {
		return err
	}
	job.SetLastTask(task)
	job.SetLastTaskPos(taskPos)
	job.SetLastTaskResult(res)

	return nil
}

func (p *V1) saveJobResult(job interfaces.Job) error {
	lastResult := job.GetLastTaskResult()
	err := p.JobResultStorage.Put(&interfaces.JobResult{
		ID:    lastResult.JobID,
		Value: lastResult.Value,
		Err:   lastResult.Err,
	})
	return err
}

func (p *V1) runRetryableJob(job interfaces.RetryableJob) {
	var tasks []interfaces.Task

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
			}
			break
		}
	}

	err := p.saveJobResult(job)
	if err != nil {
		// TODO: insert error reporting system
		return
	}
}

func (p *V1) runRevertibleJob(job interfaces.RevertibleJob) {
	var tasks []interfaces.Task
	var tasksToRevert []interfaces.RevertibleTask

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
			break
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
}

func (p *V1) runRetryableRevertibleJob(job interfaces.RetryableRevertibleJob) {
	var tasks []interfaces.Task
	var tasksToRevert []interfaces.RevertibleTask

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
			} else {
				job.SetRevertState(true)
				err := p.Queue.SetJobVisibility(job.GetID(), true)
				if err != nil {
					// TODO: insert error reporting system
					return
				}
			}
			break
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
			}
			break
		}
	}

	err := p.saveJobResult(job)
	if err != nil {
		// TODO: insert error reporting system
		return
	}
}
