//go:generate go run ../../tools/generatejob/ job.yaml quadraticequation job.gen.go

package quadraticequation

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/pkg/goodjob"
	"math"
)

// AddNumbersTask - task that takes arbitrary list of numbers and sums them.
type AddNumbersTask struct {
	ID     goodjob.TaskID
	JobID  goodjob.JobID
	Args   []*goodjob.TaskArg
	Result goodjob.TaskResult
}

func (t AddNumbersTask) GetID() goodjob.TaskID {
	return t.ID
}

func (t AddNumbersTask) GetJobID() goodjob.JobID {
	return t.JobID
}

func (t AddNumbersTask) Exec(args ...*goodjob.TaskArg) (result *goodjob.TaskResult) {
	t.Args = args
	defer func() {
		if r := recover(); r != nil {
			result = &goodjob.TaskResult{
				JobID:  t.JobID,
				TaskID: t.ID,
				Value:  nil,
				Err:    fmt.Errorf("got panic while executing the task: %v", r),
			}
		}
	}()

	var res float64

	for _, arg := range t.Args {
		res += arg.Value.(float64)
	}

	return &goodjob.TaskResult{
		JobID:  t.JobID,
		TaskID: t.ID,
		Value:  res,
		Err:    nil,
	}
}

// SubtractNumbersTask - task that takes arbitrary list of numbers and subtracts them.
type SubtractNumbersTask struct {
	ID     goodjob.TaskID
	JobID  goodjob.JobID
	Args   []*goodjob.TaskArg
	Result *goodjob.TaskResult
}

func (t SubtractNumbersTask) GetID() goodjob.TaskID {
	return t.ID
}

func (t SubtractNumbersTask) GetJobID() goodjob.JobID {
	return t.JobID
}

func (t SubtractNumbersTask) Exec(args ...*goodjob.TaskArg) (result *goodjob.TaskResult) {
	t.Args = args

	defer func() {
		if r := recover(); r != nil {
			result = &goodjob.TaskResult{
				JobID:  t.JobID,
				TaskID: t.ID,
				Value:  nil,
				Err:    fmt.Errorf("got panic while executing the task: %v", r),
			}
		}
	}()

	var res float64

	for i, arg := range t.Args {
		if i == 0 {
			res = arg.Value.(float64)
		} else {
			res -= arg.Value.(float64)
		}
	}

	return &goodjob.TaskResult{
		JobID:  t.JobID,
		TaskID: t.ID,
		Value:  res,
		Err:    nil,
	}
}

// MultiplyNumbersTask - task that takes arbitrary list of numbers and multiplies them.
type MultiplyNumbersTask struct {
	ID     goodjob.TaskID
	JobID  goodjob.JobID
	Args   []*goodjob.TaskArg
	Result *goodjob.TaskResult
}

func (t MultiplyNumbersTask) GetID() goodjob.TaskID {
	return t.ID
}

func (t MultiplyNumbersTask) GetJobID() goodjob.JobID {
	return t.JobID
}

func (t MultiplyNumbersTask) Exec(args ...*goodjob.TaskArg) (result *goodjob.TaskResult) {
	t.Args = args

	defer func() {
		if r := recover(); r != nil {
			result = &goodjob.TaskResult{
				JobID:  t.JobID,
				TaskID: t.ID,
				Value:  nil,
				Err:    fmt.Errorf("got panic while executing the task: %v", r),
			}
		}
	}()

	var res float64

	for i, arg := range t.Args {
		if i == 0 {
			res = arg.Value.(float64)
		} else {
			res *= arg.Value.(float64)
		}
	}

	return &goodjob.TaskResult{
		JobID:  t.JobID,
		TaskID: t.ID,
		Value:  res,
		Err:    nil,
	}
}

// DivideNumbersTask - task that takes arbitrary list of numbers and divides them sequentially from left to right.
type DivideNumbersTask struct {
	ID     goodjob.TaskID
	JobID  goodjob.JobID
	Args   []*goodjob.TaskArg
	Result *goodjob.TaskResult
}

func (t DivideNumbersTask) GetID() goodjob.TaskID {
	return t.ID
}

func (t DivideNumbersTask) GetJobID() goodjob.JobID {
	return t.JobID
}

func (t DivideNumbersTask) Exec(args ...*goodjob.TaskArg) (result *goodjob.TaskResult) {
	t.Args = args

	defer func() {
		if r := recover(); r != nil {
			result = &goodjob.TaskResult{
				JobID:  t.JobID,
				TaskID: t.ID,
				Value:  nil,
				Err:    fmt.Errorf("got panic while executing the task: %v", r),
			}
		}
	}()

	var res float64

	for i, arg := range t.Args {
		if arg.Value.(float64) == 0 {
			return &goodjob.TaskResult{
				JobID:  t.JobID,
				TaskID: t.ID,
				Value:  nil,
				Err:    fmt.Errorf("can't divide by zero"),
			}
		}
		if i == 0 {
			res = arg.Value.(float64)
		} else {
			res /= arg.Value.(float64)
		}
	}

	return &goodjob.TaskResult{
		JobID:  t.JobID,
		TaskID: t.ID,
		Value:  res,
		Err:    nil,
	}
}

// SquareRootOfNumberTask - task that takes arbitrary number and tries to get a square root of this number.
type SquareRootOfNumberTask struct {
	ID     goodjob.TaskID
	JobID  goodjob.JobID
	Args   []*goodjob.TaskArg
	Result *goodjob.TaskResult
}

func (t SquareRootOfNumberTask) GetID() goodjob.TaskID {
	return t.ID
}

func (t SquareRootOfNumberTask) GetJobID() goodjob.JobID {
	return t.JobID
}

func (t SquareRootOfNumberTask) Exec(args ...*goodjob.TaskArg) (result *goodjob.TaskResult) {
	t.Args = args

	defer func() {
		if r := recover(); r != nil {
			result = &goodjob.TaskResult{
				JobID:  t.JobID,
				TaskID: t.ID,
				Value:  nil,
				Err:    fmt.Errorf("got panic while executing the task: %v", r),
			}
		}
	}()

	var res float64

	res = math.Sqrt(args[0].Value.(float64))

	return &goodjob.TaskResult{
		JobID:  t.JobID,
		TaskID: t.ID,
		Value:  res,
		Err:    nil,
	}
}

// GetEquationSolutionTask - task that takes two results of previous tasks to get solution of quadratic equation.
type GetEquationSolutionTask struct {
	ID     goodjob.TaskID
	JobID  goodjob.JobID
	Args   []*goodjob.TaskArg
	Result *goodjob.TaskResult
}

func (t GetEquationSolutionTask) GetID() goodjob.TaskID {
	return t.ID
}

func (t GetEquationSolutionTask) GetJobID() goodjob.JobID {
	return t.JobID
}

func (t GetEquationSolutionTask) Exec(args ...*goodjob.TaskArg) (result *goodjob.TaskResult) {
	t.Args = args

	defer func() {
		if r := recover(); r != nil {
			result = &goodjob.TaskResult{
				JobID:  t.JobID,
				TaskID: t.ID,
				Value:  nil,
				Err:    fmt.Errorf("got panic while executing the task: %v", r),
			}
		}
	}()

	x1 := args[0].Value.(float64)
	x2 := args[1].Value.(float64)

	return &goodjob.TaskResult{
		JobID:  t.JobID,
		TaskID: t.ID,
		Value:  []float64{x1, x2},
		Err:    nil,
	}
}
