//go:generate mockgen -source ./job.go  -destination ../../mocks/job_mock.go -package mocks

package goodjob

type (
	// Job - unit of work that consist of many tasks that have to be executed in a sequential order
	Job interface {
		// GetID - return JobID
		GetID() JobID
		// GetTasks - return Job's tasks
		GetTasks() []Task
		// GetTaskArgs - return arguments that supposed to be passed to it during execution
		GetTaskArgs(taskID TaskID) []*TaskArg
		// GetVisible - is job visible for Queue processors
		GetVisible() bool
		// SetVisible - set Job visible value
		SetVisible(visible bool)
		// GetLastTask - get last run Task
		GetLastTask() Task
		// SetLastTask - set last run Task
		SetLastTask(task Task)
		// GetLastTaskPos - get position of last run Task
		GetLastTaskPos() int
		// SetLastTaskPos - set position of last run Task
		SetLastTaskPos(pos int)
		// GetLastTaskResult - get TaskResult of last run Task
		GetLastTaskResult() *TaskResult
		// SetLastTaskResult - set TaskResult of last run Task
		SetLastTaskResult(result *TaskResult)
	}

	// JobID - unique id of a Job
	JobID string

	JobArg struct {
		Name  string
		Value any
	}

	JobResult struct {
		ID    JobID
		Value any
		Err   error
	}

	// RetryableJob - a Job that can be retried from last unsuccessful Task
	RetryableJob interface {
		Job
		// GetTasksToRetry - return a list of tasks that need to be retried
		GetTasksToRetry() []Task
		// IncreaseRetryCount - increase retry attempts counter for specific Task
		IncreaseRetryCount(taskID TaskID)
		// RetryThresholdReached - is retry threshold reached for specific Task, if true retry attempts should be stopped
		RetryThresholdReached(taskID TaskID) bool
	}

	// RevertibleJob - a Job that can be reverted in case of an error encountered during Job execution
	RevertibleJob interface {
		Job
		// GetTasksToRevert - get list of task that can be reverted in case of an error encountered during Job execution
		GetTasksToRevert() []RevertibleTask
		// GetRevertState - is Job in the revert state. If true previous successful revertible tasks should be attempted to be reverted
		GetRevertState() bool
		// SetRevertState - set value for Job revert state. If set to true previous successful revertible tasks should be attempted to be reverted
		SetRevertState(revert bool)
	}

	// RetryableRevertibleJob - a Job that can be both reverted and retried in both directions of execution
	RetryableRevertibleJob interface {
		RetryableJob
		RevertibleJob
	}
)
