//go:generate mockgen -source ./queue.go  -destination ../../mocks/queue_mock.go -package mocks

package goodjob

type (
	// Queue - basic interface for Job queue.
	Queue interface {
		// AddJob - adds to Job to queue
		AddJob(job Job, args ...*QueueArg) error
		// GetNextJob - returns Job from the top of the queue.
		GetNextJob(args ...*QueueArg) (Job, error)
		// RemoveJob - removes Job from the queue
		RemoveJob(id JobID, args ...*QueueArg) error
		// SetJobVisibility - sets Job visibility inside the queue.
		// If Job visibility is set to false, consumers of the queue would not be able to get it via a call to GetNextJob
		SetJobVisibility(id JobID, visible bool, args ...*QueueArg) error
	}

	QueueArg struct {
		Name  string
		Value any
	}
)
