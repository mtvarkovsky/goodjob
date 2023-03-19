package interfaces

type (
	Processor interface {
		Start(args ...any) error
		Stop(args ...any) error
		AddJob(job Job, args ...any) error
		GetJobResult(id JobID, args ...any) (*JobResult, error)
	}
)
