// This code is automatically generated. EDIT AT YOUR OWN RISK.

package retryable_revertible

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser/dummyservices"
	"github.com/mtvarkovsky/goodjob/pkg/goodjob"
	"github.com/oklog/ulid/v2"
)

type RetryableRevertibleDeleteUserDataJob struct {
	Tasks                       []goodjob.Task
	TaskArgs                    map[goodjob.TaskID][]*goodjob.TaskArg
	Visible                     bool
	LastTask                    goodjob.Task
	LastTaskPos                 int
	LastTaskResult              *goodjob.TaskResult
	ID                          goodjob.JobID
	JobArgs                     []goodjob.JobArg
	RevertState                 bool
	ForwardRetryThreshold       int
	ForwardRetryThresholdCount  map[goodjob.TaskID]int
	BackwardRetryThreshold      int
	BackwardRetryThresholdCount map[goodjob.TaskID]int
}

func NewRetryableRevertibleDeleteUserDataJob(userID string, usersClient *dummyservices.UserServiceClient, authClient *dummyservices.AuthServiceClient, ordersClient *dummyservices.OrdersServiceClient, forwardRetryThreshold int, backwardRetryThreshold int) goodjob.RetryableRevertibleJob {
	jobID := goodjob.JobID(fmt.Sprintf("revertible delete user data with retries (%s)", ulid.Make()))
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
			Name:  "forwardRetryThreshold",
			Value: forwardRetryThreshold,
		},
		{
			Name:  "backwardRetryThreshold",
			Value: backwardRetryThreshold,
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
	return &RetryableRevertibleDeleteUserDataJob{
		JobArgs:                     jobArgs,
		LastTaskPos:                 0,
		Tasks:                       tasks,
		TaskArgs:                    taskArgs,
		Visible:                     true,
		LastTask:                    nil,
		LastTaskResult:              nil,
		ID:                          jobID,
		ForwardRetryThreshold:       forwardRetryThreshold,
		ForwardRetryThresholdCount:  make(map[goodjob.TaskID]int),
		BackwardRetryThreshold:      backwardRetryThreshold,
		BackwardRetryThresholdCount: make(map[goodjob.TaskID]int),
		RevertState:                 false,
	}
}

func (j *RetryableRevertibleDeleteUserDataJob) GetID() goodjob.JobID {
	return j.ID
}

func (j *RetryableRevertibleDeleteUserDataJob) GetTasks() []goodjob.Task {
	return j.Tasks
}

func (j *RetryableRevertibleDeleteUserDataJob) GetTaskArgs(taskID goodjob.TaskID) []*goodjob.TaskArg {
	return j.TaskArgs[taskID]
}

func (j *RetryableRevertibleDeleteUserDataJob) GetVisible() bool {
	return j.Visible
}

func (j *RetryableRevertibleDeleteUserDataJob) SetVisible(visible bool) {
	j.Visible = visible
}

func (j *RetryableRevertibleDeleteUserDataJob) GetLastTask() goodjob.Task {
	return j.LastTask
}

func (j *RetryableRevertibleDeleteUserDataJob) SetLastTask(task goodjob.Task) {
	j.LastTask = task
}

func (j *RetryableRevertibleDeleteUserDataJob) GetLastTaskPos() int {
	return j.LastTaskPos
}

func (j *RetryableRevertibleDeleteUserDataJob) SetLastTaskPos(pos int) {
	j.LastTaskPos = pos
}

func (j *RetryableRevertibleDeleteUserDataJob) GetLastTaskResult() *goodjob.TaskResult {
	return j.LastTaskResult
}

func (j *RetryableRevertibleDeleteUserDataJob) SetLastTaskResult(result *goodjob.TaskResult) {
	j.LastTaskResult = result
}

func (j *RetryableRevertibleDeleteUserDataJob) GetTasksToRetry() []goodjob.Task {
	if j.LastTaskResult.Err != nil {
		return j.Tasks[j.LastTaskPos:]
	}
	return nil
}

func (j *RetryableRevertibleDeleteUserDataJob) IncreaseRetryCount(taskID goodjob.TaskID) {
	if j.RevertState {
		j.BackwardRetryThresholdCount[taskID]++
	} else {
		j.ForwardRetryThresholdCount[taskID]++
	}
}

func (j *RetryableRevertibleDeleteUserDataJob) RetryThresholdReached(taskID goodjob.TaskID) bool {
	if j.RevertState {
		return j.BackwardRetryThresholdCount[taskID] >= j.BackwardRetryThreshold
	} else {
		return j.ForwardRetryThresholdCount[taskID] >= j.ForwardRetryThreshold
	}
}

func (j *RetryableRevertibleDeleteUserDataJob) GetTasksToRevert() []goodjob.RevertibleTask {
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

func (j *RetryableRevertibleDeleteUserDataJob) GetRevertState() bool {
	return j.RevertState
}

func (j *RetryableRevertibleDeleteUserDataJob) SetRevertState(revert bool) {
	j.RevertState = revert
}
