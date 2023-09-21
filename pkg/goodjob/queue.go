//go:generate mockgen -source ./queue.go  -destination ../../mocks/queue_mock.go -package mocks

package goodjob

type (
	// Queue - basic interface for Job queue.
	Queue interface {
		AddJob(job Job, args ...*QueueArg) error
		GetNextJob(args ...*QueueArg) (Job, error)
		RemoveJob(id JobID, args ...*QueueArg) error
		SetJobVisibility(id JobID, visible bool, args ...*QueueArg) error
	}

	QueueArg struct {
		Name  string
		Value any
	}
)
