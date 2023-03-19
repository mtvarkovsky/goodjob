//go:generate mockgen -source ./queue.go  -destination ../../mocks/queue_mock.go -package mocks

package interfaces

type (
	// Queue - basic interface for Job queue.
	Queue interface {
		AddJob(job Job, args ...any) error
		GetNextJob(args ...any) (Job, error)
		RemoveJob(id JobID, args ...any) error
		SetJobVisibility(id JobID, visible bool, args ...any) error
	}
)
