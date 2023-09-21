//go:generate mockgen -source ./task.go  -destination ../../mocks/task_mock.go -package mocks

package goodjob

type (
	// Task - an atomic unit of work within a job
	Task interface {
		// GetID - get TaskID of a task
		GetID() TaskID
		// GetJobID - get JobID of a Job that task belongs to
		GetJobID() JobID
		// Exec - execute task with optional arguments
		Exec(args ...*TaskArg) *TaskResult
	}

	// TaskID - id of a task
	TaskID string

	// TaskArg - argument for a Task
	TaskArg struct {
		// Name - argument name
		Name string
		// Value - argument value
		Value any
		// ValueFrom - TaskID of a task which result is supposed to be used as a value
		ValueFrom TaskID
	}

	// TaskResult - result of Task execution
	TaskResult struct {
		// JobID - id of a parent Job
		JobID JobID
		// TaskID - id of a Task
		TaskID TaskID
		// Value - actual task result
		Value any
		// Err - error that have occurred during Task execution
		Err error
	}

	// RevertibleTask - a task that can be reverted
	RevertibleTask interface {
		Task
		// Revert - revert results of Task execution
		Revert(args ...*TaskArg) (result *TaskResult)
	}
)
