// This code is automatically generated. EDIT AT YOUR OWN RISK.

package simple

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser/dummyservices"
	"github.com/mtvarkovsky/goodjob/pkg/goodjob"
	"github.com/oklog/ulid/v2"
)

type SimpleDeleteUserDataJob struct {
	ID             goodjob.JobID
	JobArgs        []goodjob.JobArg
	Tasks          []goodjob.Task
	TaskArgs       map[goodjob.TaskID][]*goodjob.TaskArg
	Visible        bool
	LastTask       goodjob.Task
	LastTaskPos    int
	LastTaskResult *goodjob.TaskResult
}

func NewSimpleDeleteUserDataJob(userID string, usersClient *dummyservices.UserServiceClient, authClient *dummyservices.AuthServiceClient, ordersClient *dummyservices.OrdersServiceClient) goodjob.Job {
	jobID := goodjob.JobID(fmt.Sprintf("delete user data (%s)", ulid.Make()))
	jobArgs := []goodjob.JobArg{
		{
			Name:  "userID",
			Value: userID,
		},
		{
			Name:  "usersClient",
			Value: usersClient,
		},
		{
			Name:  "authClient",
			Value: authClient,
		},
		{
			Name:  "ordersClient",
			Value: ordersClient,
		},
	}

	taskIDs := map[string]goodjob.TaskID{
		"safe delete user auth data":   goodjob.TaskID(fmt.Sprintf("safe delete user auth data (%s)", ulid.Make())),
		"safe delete user data":        goodjob.TaskID(fmt.Sprintf("safe delete user data (%s)", ulid.Make())),
		"safe delete user orders data": goodjob.TaskID(fmt.Sprintf("safe delete user orders data (%s)", ulid.Make())),
	}
	tasks := []goodjob.Task{
		deleteuser.SafeDeleteAuthDataTask{
			ID:    taskIDs["safe delete user auth data"],
			JobID: jobID,
		},
		deleteuser.SafeDeleteUserTask{
			ID:    taskIDs["safe delete user data"],
			JobID: jobID,
		},
		deleteuser.SafeDeleteUserOrdersTask{
			ID:    taskIDs["safe delete user orders data"],
			JobID: jobID,
		},
	}
	taskArgs := map[goodjob.TaskID][]*goodjob.TaskArg{
		taskIDs["safe delete user auth data"]: {
			{
				Name:  "authClient",
				Value: authClient,
			},
			{
				Name:  "userID",
				Value: userID,
			},
		},
		taskIDs["safe delete user data"]: {
			{
				Name:  "usersClient",
				Value: usersClient,
			},
			{
				Name:  "userID",
				Value: userID,
			},
		},
		taskIDs["safe delete user orders data"]: {
			{
				Name:  "ordersClient",
				Value: ordersClient,
			},
			{
				Name:  "userID",
				Value: userID,
			},
		},
	}
	return &SimpleDeleteUserDataJob{
		JobArgs:        jobArgs,
		LastTaskPos:    0,
		Tasks:          tasks,
		TaskArgs:       taskArgs,
		Visible:        true,
		LastTask:       nil,
		LastTaskResult: nil,
		ID:             jobID,
	}
}

func (j *SimpleDeleteUserDataJob) GetID() goodjob.JobID {
	return j.ID
}

func (j *SimpleDeleteUserDataJob) GetTasks() []goodjob.Task {
	return j.Tasks
}

func (j *SimpleDeleteUserDataJob) GetTaskArgs(taskID goodjob.TaskID) []*goodjob.TaskArg {
	return j.TaskArgs[taskID]
}

func (j *SimpleDeleteUserDataJob) GetVisible() bool {
	return j.Visible
}

func (j *SimpleDeleteUserDataJob) SetVisible(visible bool) {
	j.Visible = visible
}

func (j *SimpleDeleteUserDataJob) GetLastTask() goodjob.Task {
	return j.LastTask
}

func (j *SimpleDeleteUserDataJob) SetLastTask(task goodjob.Task) {
	j.LastTask = task
}

func (j *SimpleDeleteUserDataJob) GetLastTaskPos() int {
	return j.LastTaskPos
}

func (j *SimpleDeleteUserDataJob) SetLastTaskPos(pos int) {
	j.LastTaskPos = pos
}

func (j *SimpleDeleteUserDataJob) GetLastTaskResult() *goodjob.TaskResult {
	return j.LastTaskResult
}

func (j *SimpleDeleteUserDataJob) SetLastTaskResult(result *goodjob.TaskResult) {
	j.LastTaskResult = result
}
