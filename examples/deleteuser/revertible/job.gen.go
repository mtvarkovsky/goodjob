// This code is automatically generated. EDIT AT YOUR OWN RISK.

package revertible

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser/dummyservices"
	"github.com/mtvarkovsky/goodjob/pkg/interfaces"
	"github.com/oklog/ulid/v2"
)

type RevertibleDeleteUserDataJob struct {
	Tasks          []interfaces.Task
	TaskArgs       map[interfaces.TaskID][]*interfaces.TaskArg
	Visible        bool
	LastTask       interfaces.Task
	LastTaskPos    int
	LastTaskResult *interfaces.TaskResult
	ID             interfaces.JobID
	JobArgs        []interfaces.JobArg
	RevertState    bool
}

func NewRevertibleDeleteUserDataJob(userID string, usersClient *dummyservices.UserServiceClient, authClient *dummyservices.AuthServiceClient, ordersClient *dummyservices.OrdersServiceClient) interfaces.RevertibleJob {
	jobID := interfaces.JobID(fmt.Sprintf("delete user data with revert on failure (%s)", ulid.Make()))
	jobArgs := []interfaces.JobArg{
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

	taskIDs := map[string]interfaces.TaskID{
		"safe delete user auth data":   interfaces.TaskID(fmt.Sprintf("safe delete user auth data (%s)", ulid.Make())),
		"safe delete user data":        interfaces.TaskID(fmt.Sprintf("safe delete user data (%s)", ulid.Make())),
		"safe delete user orders data": interfaces.TaskID(fmt.Sprintf("safe delete user orders data (%s)", ulid.Make())),
	}
	tasks := []interfaces.Task{
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
	taskArgs := map[interfaces.TaskID][]*interfaces.TaskArg{
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
	return &RevertibleDeleteUserDataJob{
		JobArgs:        jobArgs,
		LastTaskPos:    0,
		Tasks:          tasks,
		TaskArgs:       taskArgs,
		Visible:        true,
		LastTask:       nil,
		LastTaskResult: nil,
		ID:             jobID,
		RevertState:    false,
	}
}

func (j *RevertibleDeleteUserDataJob) GetID() interfaces.JobID {
	return j.ID
}

func (j *RevertibleDeleteUserDataJob) GetTasks() []interfaces.Task {
	return j.Tasks
}

func (j *RevertibleDeleteUserDataJob) GetTaskArgs(taskID interfaces.TaskID) []*interfaces.TaskArg {
	return j.TaskArgs[taskID]
}

func (j *RevertibleDeleteUserDataJob) GetVisible() bool {
	return j.Visible
}

func (j *RevertibleDeleteUserDataJob) SetVisible(visible bool) {
	j.Visible = visible
}

func (j *RevertibleDeleteUserDataJob) GetLastTask() interfaces.Task {
	return j.LastTask
}

func (j *RevertibleDeleteUserDataJob) SetLastTask(task interfaces.Task) {
	j.LastTask = task
}

func (j *RevertibleDeleteUserDataJob) GetLastTaskPos() int {
	return j.LastTaskPos
}

func (j *RevertibleDeleteUserDataJob) SetLastTaskPos(pos int) {
	j.LastTaskPos = pos
}

func (j *RevertibleDeleteUserDataJob) GetLastTaskResult() *interfaces.TaskResult {
	return j.LastTaskResult
}

func (j *RevertibleDeleteUserDataJob) SetLastTaskResult(result *interfaces.TaskResult) {
	j.LastTaskResult = result
}

func (j *RevertibleDeleteUserDataJob) GetTasksToRevert() []interfaces.RevertibleTask {
	var tasksToRevert []interfaces.RevertibleTask
	for i := j.LastTaskPos - 1; i >= 0; i-- {
		task := j.Tasks[i]
		switch t := task.(type) {
		case interfaces.RevertibleTask:
			tasksToRevert = append(tasksToRevert, t)
		}
	}
	return tasksToRevert
}

func (j *RevertibleDeleteUserDataJob) GetRevertState() bool {
	return j.RevertState
}

func (j *RevertibleDeleteUserDataJob) SetRevertState(revert bool) {
	j.RevertState = revert
}
