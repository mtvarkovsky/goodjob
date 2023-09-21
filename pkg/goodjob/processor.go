package goodjob

type (
	Processor interface {
		Start(args ...*ProcessorArg) error
		Stop(args ...*ProcessorArg) error
		AddJob(job Job, args ...*ProcessorArg) error
		GetJobResult(id JobID, args ...*ProcessorArg) (*JobResult, error)
	}

	ProcessorArg struct {
		Name  string
		Value any
	}
)
