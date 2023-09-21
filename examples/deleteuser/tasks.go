//go:generate go run ../../tools/generatejob/ ./simple/job.yaml simple ./simple/job.gen.go
//go:generate go run ../../tools/generatejob/ ./retryable/job.yaml retryable ./retryable/job.gen.go
//go:generate go run ../../tools/generatejob/ ./revertible/job.yaml revertible ./revertible/job.gen.go
//go:generate go run ../../tools/generatejob/ ./retryable_revertible/job.yaml retryable_revertible ./retryable_revertible/job.gen.go

package deleteuser

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser/dummyservices"
	"github.com/mtvarkovsky/goodjob/pkg/goodjob"
)

type SafeDeleteUserTask struct {
	ID    goodjob.TaskID
	JobID goodjob.JobID
	Args  []*goodjob.TaskArg
}

func (t SafeDeleteUserTask) GetID() goodjob.TaskID {
	return t.ID
}

func (t SafeDeleteUserTask) GetJobID() goodjob.JobID {
	return t.JobID
}

func (t SafeDeleteUserTask) Exec(args ...*goodjob.TaskArg) (result *goodjob.TaskResult) {
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

	client := args[0].Value.(*dummyservices.UserServiceClient)
	userID := args[1].Value.(string)

	request := &dummyservices.SafeDeleteUserRequest{
		UserID: userID,
	}

	err := client.SafeDeleteUser(request)

	return &goodjob.TaskResult{
		TaskID: t.ID,
		JobID:  t.JobID,
		Value:  nil,
		Err:    err,
	}
}

func (t SafeDeleteUserTask) Revert(args ...*goodjob.TaskArg) (result *goodjob.TaskResult) {
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

	client := args[0].Value.(*dummyservices.UserServiceClient)
	userID := args[1].Value.(string)

	request := &dummyservices.RestoreUserRequest{
		UserID: userID,
	}

	err := client.RestoreUser(request)

	return &goodjob.TaskResult{
		TaskID: t.ID,
		JobID:  t.JobID,
		Value:  nil,
		Err:    err,
	}
}

type SafeDeleteAuthDataTask struct {
	ID    goodjob.TaskID
	JobID goodjob.JobID
	Args  []*goodjob.TaskArg
}

func (t SafeDeleteAuthDataTask) GetID() goodjob.TaskID {
	return t.ID
}

func (t SafeDeleteAuthDataTask) GetJobID() goodjob.JobID {
	return t.JobID
}

func (t SafeDeleteAuthDataTask) Exec(args ...*goodjob.TaskArg) (result *goodjob.TaskResult) {
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

	client := args[0].Value.(*dummyservices.AuthServiceClient)
	userID := args[1].Value.(string)

	request := &dummyservices.SafeDeleteAuthDataRequest{
		UserID: userID,
	}

	err := client.SafeDeleteAuthData(request)

	return &goodjob.TaskResult{
		TaskID: t.ID,
		JobID:  t.JobID,
		Value:  nil,
		Err:    err,
	}
}

func (t SafeDeleteAuthDataTask) Revert(args ...*goodjob.TaskArg) (result *goodjob.TaskResult) {
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

	client := args[0].Value.(*dummyservices.AuthServiceClient)
	userID := args[1].Value.(string)

	request := &dummyservices.RestoreAuthDataRequest{
		UserID: userID,
	}

	err := client.RestoreAuthData(request)

	return &goodjob.TaskResult{
		TaskID: t.ID,
		JobID:  t.JobID,
		Value:  nil,
		Err:    err,
	}
}

type SafeDeleteUserOrdersTask struct {
	ID    goodjob.TaskID
	JobID goodjob.JobID
	Args  []*goodjob.TaskArg
}

func (t SafeDeleteUserOrdersTask) GetID() goodjob.TaskID {
	return t.ID
}

func (t SafeDeleteUserOrdersTask) GetJobID() goodjob.JobID {
	return t.JobID
}

func (t SafeDeleteUserOrdersTask) Exec(args ...*goodjob.TaskArg) (result *goodjob.TaskResult) {
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

	client := args[0].Value.(*dummyservices.OrdersServiceClient)
	userID := args[1].Value.(string)

	request := &dummyservices.SafeDeleteUserOrdersRequest{
		UserID: userID,
	}

	err := client.SafeDeleteUserOrders(request)

	return &goodjob.TaskResult{
		TaskID: t.ID,
		JobID:  t.JobID,
		Value:  nil,
		Err:    err,
	}
}
