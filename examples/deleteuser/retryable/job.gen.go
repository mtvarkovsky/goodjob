// This code is automatically generated. EDIT AT YOUR OWN RISK.

package retryable

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser/dummyservices"
	"github.com/mtvarkovsky/goodjob/pkg/interfaces"
	"github.com/oklog/ulid/v2"
)

type RetryableDeleteUserDataJob struct {
	ID                  interfaces.JobID
	JobArgs             []interfaces.JobArg
	Tasks               []interfaces.Task
	TaskArgs            map[interfaces.TaskID][]*interfaces.TaskArg
	Visible             bool
	LastTask            interfaces.Task
	LastTaskPos         int
	LastTaskResult      *interfaces.TaskResult
	RetryThreshold      int
	RetryThresholdCount map[interfaces.TaskID]int
}

func NewRetryableDeleteUserDataJob(userID string, usersClient *dummyservices.UserServiceClient, authClient *dummyservices.AuthServiceClient, ordersClient *dummyservices.OrdersServiceClient, retryThreshold int) interfaces.RetryableJob {
	jobID := interfaces.JobID(fmt.Sprintf("delete user data with retries (%s)", ulid.Make()))
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
		{
			Name:  "retryThreshold",
			Value: retryThreshold,
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
	return &RetryableDeleteUserDataJob{
		JobArgs:             jobArgs,
		LastTaskPos:         0,
		Tasks:               tasks,
		TaskArgs:            taskArgs,
		Visible:             true,
		LastTask:            nil,
		LastTaskResult:      nil,
		ID:                  jobID,
		RetryThreshold:      retryThreshold,
		RetryThresholdCount: make(map[interfaces.TaskID]int),
	}
}

func (j *RetryableDeleteUserDataJob) GetID() interfaces.JobID {
	return j.ID
}

func (j *RetryableDeleteUserDataJob) GetTasks() []interfaces.Task {
	return j.Tasks
}

func (j *RetryableDeleteUserDataJob) GetTaskArgs(taskID interfaces.TaskID) []*interfaces.TaskArg {
	return j.TaskArgs[taskID]
}

func (j *RetryableDeleteUserDataJob) GetVisible() bool {
	return j.Visible
}

func (j *RetryableDeleteUserDataJob) SetVisible(visible bool) {
	j.Visible = visible
}

func (j *RetryableDeleteUserDataJob) GetLastTask() interfaces.Task {
	return j.LastTask
}

func (j *RetryableDeleteUserDataJob) SetLastTask(task interfaces.Task) {
	j.LastTask = task
}

func (j *RetryableDeleteUserDataJob) GetLastTaskPos() int {
	return j.LastTaskPos
}

func (j *RetryableDeleteUserDataJob) SetLastTaskPos(pos int) {
	j.LastTaskPos = pos
}

func (j *RetryableDeleteUserDataJob) GetLastTaskResult() *interfaces.TaskResult {
	return j.LastTaskResult
}

func (j *RetryableDeleteUserDataJob) SetLastTaskResult(result *interfaces.TaskResult) {
	j.LastTaskResult = result
}

func (j *RetryableDeleteUserDataJob) GetTasksToRetry() []interfaces.Task {
	if j.LastTaskResult.Err != nil {
		return j.Tasks[j.LastTaskPos:]
	}
	return nil
}

func (j *RetryableDeleteUserDataJob) IncreaseRetryCount(taskID interfaces.TaskID) {
	j.RetryThresholdCount[taskID]++
}

func (j *RetryableDeleteUserDataJob) RetryThresholdReached(taskID interfaces.TaskID) bool {
	return j.RetryThresholdCount[taskID] >= j.RetryThreshold
}
