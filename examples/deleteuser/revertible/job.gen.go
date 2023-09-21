// This code is automatically generated. EDIT AT YOUR OWN RISK.

package revertible

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser/dummyservices"
	"github.com/mtvarkovsky/goodjob/pkg/goodjob"
	"github.com/oklog/ulid/v2"
)

type RevertibleDeleteUserDataJob struct {
	Tasks          []goodjob.Task
	TaskArgs       map[goodjob.TaskID][]*goodjob.TaskArg
	Visible        bool
	LastTask       goodjob.Task
	LastTaskPos    int
	LastTaskResult *goodjob.TaskResult
	ID             goodjob.JobID
	JobArgs        []goodjob.JobArg
	RevertState    bool
}

func NewRevertibleDeleteUserDataJob(userID string, usersClient *dummyservices.UserServiceClient, authClient *dummyservices.AuthServiceClient, ordersClient *dummyservices.OrdersServiceClient) goodjob.RevertibleJob {
	jobID := goodjob.JobID(fmt.Sprintf("delete user data with revert on failure (%s)", ulid.Make()))
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

func (j *RevertibleDeleteUserDataJob) GetID() goodjob.JobID {
	return j.ID
}

func (j *RevertibleDeleteUserDataJob) GetTasks() []goodjob.Task {
	return j.Tasks
}

func (j *RevertibleDeleteUserDataJob) GetTaskArgs(taskID goodjob.TaskID) []*goodjob.TaskArg {
	return j.TaskArgs[taskID]
}

func (j *RevertibleDeleteUserDataJob) GetVisible() bool {
	return j.Visible
}

func (j *RevertibleDeleteUserDataJob) SetVisible(visible bool) {
	j.Visible = visible
}

func (j *RevertibleDeleteUserDataJob) GetLastTask() goodjob.Task {
	return j.LastTask
}

func (j *RevertibleDeleteUserDataJob) SetLastTask(task goodjob.Task) {
	j.LastTask = task
}

func (j *RevertibleDeleteUserDataJob) GetLastTaskPos() int {
	return j.LastTaskPos
}

func (j *RevertibleDeleteUserDataJob) SetLastTaskPos(pos int) {
	j.LastTaskPos = pos
}

func (j *RevertibleDeleteUserDataJob) GetLastTaskResult() *goodjob.TaskResult {
	return j.LastTaskResult
}

func (j *RevertibleDeleteUserDataJob) SetLastTaskResult(result *goodjob.TaskResult) {
	j.LastTaskResult = result
}

func (j *RevertibleDeleteUserDataJob) GetTasksToRevert() []goodjob.RevertibleTask {
	var tasksToRevert []goodjob.RevertibleTask
	for i := j.LastTaskPos - 1; i >= 0; i-- {
		task := j.Tasks[i]
		switch t := task.(type) {
		case goodjob.RevertibleTask:
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
