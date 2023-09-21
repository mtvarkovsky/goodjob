//go:generate mockgen -source ./storage.go  -destination ../../mocks/storage_mock.go -package mocks

package goodjob

type (
	TaskResultsStorage interface {
		Put(result *TaskResult) error
		Get(jobID JobID, taskID TaskID) (*TaskResult, error)
	}

	JobResultsStorage interface {
		Put(result *JobResult) error
		Get(id JobID) (*JobResult, error)
	}
)
