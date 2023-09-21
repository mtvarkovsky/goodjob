// This code is automatically generated. EDIT AT YOUR OWN RISK.

package retryable

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser/dummyservices"
	"github.com/mtvarkovsky/goodjob/pkg/goodjob"
	"github.com/oklog/ulid/v2"
)

type RetryableDeleteUserDataJob struct {
	LastTaskPos         int
	LastTaskResult      *goodjob.TaskResult
	ID                  goodjob.JobID
	JobArgs             []goodjob.JobArg
	Tasks               []goodjob.Task
	TaskArgs            map[goodjob.TaskID][]*goodjob.TaskArg
	Visible             bool
	LastTask            goodjob.Task
	RetryThreshold      int
	RetryThresholdCount map[goodjob.TaskID]int
}

func NewRetryableDeleteUserDataJob(userID string, usersClient *dummyservices.UserServiceClient, authClient *dummyservices.AuthServiceClient, ordersClient *dummyservices.OrdersServiceClient, retryThreshold int) goodjob.RetryableJob {
	jobID := goodjob.JobID(fmt.Sprintf("delete user data with retries (%s)", ulid.Make()))
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
		{
			Name:  "retryThreshold",
			Value: retryThreshold,
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
		RetryThresholdCount: make(map[goodjob.TaskID]int),
	}
}

func (j *RetryableDeleteUserDataJob) GetID() goodjob.JobID {
	return j.ID
}

func (j *RetryableDeleteUserDataJob) GetTasks() []goodjob.Task {
	return j.Tasks
}

func (j *RetryableDeleteUserDataJob) GetTaskArgs(taskID goodjob.TaskID) []*goodjob.TaskArg {
	return j.TaskArgs[taskID]
}

func (j *RetryableDeleteUserDataJob) GetVisible() bool {
	return j.Visible
}

func (j *RetryableDeleteUserDataJob) SetVisible(visible bool) {
	j.Visible = visible
}

func (j *RetryableDeleteUserDataJob) GetLastTask() goodjob.Task {
	return j.LastTask
}

func (j *RetryableDeleteUserDataJob) SetLastTask(task goodjob.Task) {
	j.LastTask = task
}

func (j *RetryableDeleteUserDataJob) GetLastTaskPos() int {
	return j.LastTaskPos
}

func (j *RetryableDeleteUserDataJob) SetLastTaskPos(pos int) {
	j.LastTaskPos = pos
}

func (j *RetryableDeleteUserDataJob) GetLastTaskResult() *goodjob.TaskResult {
	return j.LastTaskResult
}

func (j *RetryableDeleteUserDataJob) SetLastTaskResult(result *goodjob.TaskResult) {
	j.LastTaskResult = result
}

func (j *RetryableDeleteUserDataJob) GetTasksToRetry() []goodjob.Task {
	if j.LastTaskResult.Err != nil {
		return j.Tasks[j.LastTaskPos:]
	}
	return nil
}

func (j *RetryableDeleteUserDataJob) IncreaseRetryCount(taskID goodjob.TaskID) {
	j.RetryThresholdCount[taskID]++
}

func (j *RetryableDeleteUserDataJob) RetryThresholdReached(taskID goodjob.TaskID) bool {
	return j.RetryThresholdCount[taskID] >= j.RetryThreshold
}
