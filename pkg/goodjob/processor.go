package goodjob

type (
	// Processor - is a worker that accepts jobs for execution, executes them and allows to get result of executed jobs.
	Processor interface {
		// Start - starts the processor
		Start(args ...*ProcessorArg) error
		// Stop - stops the processor
		Stop(args ...*ProcessorArg) error
		// AddJob - adds job to be processed
		AddJob(job Job, args ...*ProcessorArg) error
		// GetJobResult - retrieves job result. If job result not found, an error is returned
		GetJobResult(id JobID, args ...*ProcessorArg) (*JobResult, error)
	}

	ProcessorArg struct {
		Name  string
		Value any
	}
)
