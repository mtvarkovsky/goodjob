//go:generate go run ../../tools/generatejob/ job.yaml quadraticequation job.gen.go

package quadraticequation

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/pkg/interfaces"
	"math"
)

// AddNumbersTask - task that takes arbitrary list of numbers and sums them.
type AddNumbersTask struct {
	ID     interfaces.TaskID
	JobID  interfaces.JobID
	Args   []*interfaces.TaskArg
	Result interfaces.TaskResult
}

func (t AddNumbersTask) GetID() interfaces.TaskID {
	return t.ID
}

func (t AddNumbersTask) GetJobID() interfaces.JobID {
	return t.JobID
}

func (t AddNumbersTask) Exec(args ...*interfaces.TaskArg) (result *interfaces.TaskResult) {
	t.Args = args
	defer func() {
		if r := recover(); r != nil {
			result = &interfaces.TaskResult{
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

	return &interfaces.TaskResult{
		JobID:  t.JobID,
		TaskID: t.ID,
		Value:  res,
		Err:    nil,
	}
}

// SubtractNumbersTask - task that takes arbitrary list of numbers and subtracts them.
type SubtractNumbersTask struct {
	ID     interfaces.TaskID
	JobID  interfaces.JobID
	Args   []*interfaces.TaskArg
	Result *interfaces.TaskResult
}

func (t SubtractNumbersTask) GetID() interfaces.TaskID {
	return t.ID
}

func (t SubtractNumbersTask) GetJobID() interfaces.JobID {
	return t.JobID
}

func (t SubtractNumbersTask) Exec(args ...*interfaces.TaskArg) (result *interfaces.TaskResult) {
	t.Args = args

	defer func() {
		if r := recover(); r != nil {
			result = &interfaces.TaskResult{
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

	return &interfaces.TaskResult{
		JobID:  t.JobID,
		TaskID: t.ID,
		Value:  res,
		Err:    nil,
	}
}

// MultiplyNumbersTask - task that takes arbitrary list of numbers and multiplies them.
type MultiplyNumbersTask struct {
	ID     interfaces.TaskID
	JobID  interfaces.JobID
	Args   []*interfaces.TaskArg
	Result *interfaces.TaskResult
}

func (t MultiplyNumbersTask) GetID() interfaces.TaskID {
	return t.ID
}

func (t MultiplyNumbersTask) GetJobID() interfaces.JobID {
	return t.JobID
}

func (t MultiplyNumbersTask) Exec(args ...*interfaces.TaskArg) (result *interfaces.TaskResult) {
	t.Args = args

	defer func() {
		if r := recover(); r != nil {
			result = &interfaces.TaskResult{
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

	return &interfaces.TaskResult{
		JobID:  t.JobID,
		TaskID: t.ID,
		Value:  res,
		Err:    nil,
	}
}

// DivideNumbersTask - task that takes arbitrary list of numbers and divides them sequentially from left to right.
type DivideNumbersTask struct {
	ID     interfaces.TaskID
	JobID  interfaces.JobID
	Args   []*interfaces.TaskArg
	Result *interfaces.TaskResult
}

func (t DivideNumbersTask) GetID() interfaces.TaskID {
	return t.ID
}

func (t DivideNumbersTask) GetJobID() interfaces.JobID {
	return t.JobID
}

func (t DivideNumbersTask) Exec(args ...*interfaces.TaskArg) (result *interfaces.TaskResult) {
	t.Args = args

	defer func() {
		if r := recover(); r != nil {
			result = &interfaces.TaskResult{
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
			return &interfaces.TaskResult{
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

	return &interfaces.TaskResult{
		JobID:  t.JobID,
		TaskID: t.ID,
		Value:  res,
		Err:    nil,
	}
}

// SquareRootOfNumberTask - task that takes arbitrary number and tries to get a square root of this number.
type SquareRootOfNumberTask struct {
	ID     interfaces.TaskID
	JobID  interfaces.JobID
	Args   []*interfaces.TaskArg
	Result *interfaces.TaskResult
}

func (t SquareRootOfNumberTask) GetID() interfaces.TaskID {
	return t.ID
}

func (t SquareRootOfNumberTask) GetJobID() interfaces.JobID {
	return t.JobID
}

func (t SquareRootOfNumberTask) Exec(args ...*interfaces.TaskArg) (result *interfaces.TaskResult) {
	t.Args = args

	defer func() {
		if r := recover(); r != nil {
			result = &interfaces.TaskResult{
				JobID:  t.JobID,
				TaskID: t.ID,
				Value:  nil,
				Err:    fmt.Errorf("got panic while executing the task: %v", r),
			}
		}
	}()

	var res float64

	res = math.Sqrt(args[0].Value.(float64))

	return &interfaces.TaskResult{
		JobID:  t.JobID,
		TaskID: t.ID,
		Value:  res,
		Err:    nil,
	}
}

// GetEquationSolutionTask - task that takes two results of previous tasks to get solution of quadratic equation.
type GetEquationSolutionTask struct {
	ID     interfaces.TaskID
	JobID  interfaces.JobID
	Args   []*interfaces.TaskArg
	Result *interfaces.TaskResult
}

func (t GetEquationSolutionTask) GetID() interfaces.TaskID {
	return t.ID
}

func (t GetEquationSolutionTask) GetJobID() interfaces.JobID {
	return t.JobID
}

func (t GetEquationSolutionTask) Exec(args ...*interfaces.TaskArg) (result *interfaces.TaskResult) {
	t.Args = args

	defer func() {
		if r := recover(); r != nil {
			result = &interfaces.TaskResult{
				JobID:  t.JobID,
				TaskID: t.ID,
				Value:  nil,
				Err:    fmt.Errorf("got panic while executing the task: %v", r),
			}
		}
	}()

	x1 := args[0].Value.(float64)
	x2 := args[1].Value.(float64)

	return &interfaces.TaskResult{
		JobID:  t.JobID,
		TaskID: t.ID,
		Value:  []float64{x1, x2},
		Err:    nil,
	}
}
