//go:generate go run ../../tools/generatejob/ ./simple/job.yaml simple ./simple/job.gen.go
//go:generate go run ../../tools/generatejob/ ./retryable/job.yaml retryable ./retryable/job.gen.go
//go:generate go run ../../tools/generatejob/ ./revertible/job.yaml revertible ./revertible/job.gen.go
//go:generate go run ../../tools/generatejob/ ./retryable_revertible/job.yaml retryable_revertible ./retryable_revertible/job.gen.go


package deleteuser

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser/dummyservices"
	"github.com/mtvarkovsky/goodjob/pkg/interfaces"
)

type SafeDeleteUserTask struct {
	ID    interfaces.TaskID
	JobID interfaces.JobID
	Args  []*interfaces.TaskArg
}

func (t SafeDeleteUserTask) GetID() interfaces.TaskID {
	return t.ID
}

func (t SafeDeleteUserTask) GetJobID() interfaces.JobID {
	return t.JobID
}

func (t SafeDeleteUserTask) Exec(args ...*interfaces.TaskArg) (result *interfaces.TaskResult) {
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

	client := args[0].Value.(*dummyservices.UserServiceClient)
	userID := args[1].Value.(string)

	request := &dummyservices.SafeDeleteUserRequest{
		UserID: userID,
	}

	err := client.SafeDeleteUser(request)

	return &interfaces.TaskResult{
		TaskID: t.ID,
		JobID:  t.JobID,
		Value:  nil,
		Err:    err,
	}
}

func (t SafeDeleteUserTask) Revert(args ...*interfaces.TaskArg) (result *interfaces.TaskResult) {
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

	client := args[0].Value.(*dummyservices.UserServiceClient)
	userID := args[1].Value.(string)

	request := &dummyservices.RestoreUserRequest{
		UserID: userID,
	}

	err := client.RestoreUser(request)

	return &interfaces.TaskResult{
		TaskID: t.ID,
		JobID:  t.JobID,
		Value:  nil,
		Err:    err,
	}
}

type SafeDeleteAuthDataTask struct {
	ID    interfaces.TaskID
	JobID interfaces.JobID
	Args  []*interfaces.TaskArg
}

func (t SafeDeleteAuthDataTask) GetID() interfaces.TaskID {
	return t.ID
}

func (t SafeDeleteAuthDataTask) GetJobID() interfaces.JobID {
	return t.JobID
}

func (t SafeDeleteAuthDataTask) Exec(args ...*interfaces.TaskArg) (result *interfaces.TaskResult) {
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

	client := args[0].Value.(*dummyservices.AuthServiceClient)
	userID := args[1].Value.(string)

	request := &dummyservices.SafeDeleteAuthDataRequest{
		UserID: userID,
	}

	err := client.SafeDeleteAuthData(request)

	return &interfaces.TaskResult{
		TaskID: t.ID,
		JobID:  t.JobID,
		Value:  nil,
		Err:    err,
	}
}

func (t SafeDeleteAuthDataTask) Revert(args ...*interfaces.TaskArg) (result *interfaces.TaskResult) {
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

	client := args[0].Value.(*dummyservices.AuthServiceClient)
	userID := args[1].Value.(string)

	request := &dummyservices.RestoreAuthDataRequest{
		UserID: userID,
	}

	err := client.RestoreAuthData(request)

	return &interfaces.TaskResult{
		TaskID: t.ID,
		JobID:  t.JobID,
		Value:  nil,
		Err:    err,
	}
}

type SafeDeleteUserOrdersTask struct {
	ID    interfaces.TaskID
	JobID interfaces.JobID
	Args  []*interfaces.TaskArg
}

func (t SafeDeleteUserOrdersTask) GetID() interfaces.TaskID {
	return t.ID
}

func (t SafeDeleteUserOrdersTask) GetJobID() interfaces.JobID {
	return t.JobID
}

func (t SafeDeleteUserOrdersTask) Exec(args ...*interfaces.TaskArg) (result *interfaces.TaskResult) {
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

	client := args[0].Value.(*dummyservices.OrdersServiceClient)
	userID := args[1].Value.(string)

	request := &dummyservices.SafeDeleteUserOrdersRequest{
		UserID: userID,
	}

	err := client.SafeDeleteUserOrders(request)

	return &interfaces.TaskResult{
		TaskID: t.ID,
		JobID:  t.JobID,
		Value:  nil,
		Err:    err,
	}
}
