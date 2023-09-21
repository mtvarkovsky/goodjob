package retryable_revertible

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/examples/deleteuser/dummyservices"
	"github.com/mtvarkovsky/goodjob/pkg/goodjob"
	"github.com/mtvarkovsky/goodjob/pkg/processor"
	"github.com/mtvarkovsky/goodjob/pkg/queue"
	"github.com/mtvarkovsky/goodjob/pkg/storage"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TODO: rework tests to remove time.Sleep()

func TestRetryableRevertibleDeleteUser_Success(t *testing.T) {
	userID := "user123"

	forwardRetryCountOrders := 4
	backwardRetryCountUsers := 3
	backwardRetryCountAuthData := 3

	var taskActionsLog []string

	usersClient := dummyservices.NewUserServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call safe delete user")
			return nil
		},
		func() error {
			if backwardRetryCountUsers > 0 {
				taskActionsLog = append(taskActionsLog, "failed call revert user data")
				backwardRetryCountUsers--
				return fmt.Errorf("cant revert user data")
			} else {
				taskActionsLog = append(taskActionsLog, "successfully call revert user data")
				return nil
			}
		},
	)

	authClient := dummyservices.NewAuthServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call safe delete user auth data")
			return nil
		},
		func() error {
			if backwardRetryCountAuthData > 0 {
				taskActionsLog = append(taskActionsLog, "failed call revert user auth data")
				backwardRetryCountAuthData--
				return fmt.Errorf("cant revert user auth data")
			} else {
				taskActionsLog = append(taskActionsLog, "successfully call revert user auth data")
				return nil
			}
		},
	)

	ordersClient := dummyservices.NewOrdersServiceClient(
		func() error {
			if forwardRetryCountOrders > 0 {
				taskActionsLog = append(taskActionsLog, "failed call safe delete user orders data")
				forwardRetryCountOrders--
				return fmt.Errorf("cant delete user orders data")
			} else {
				taskActionsLog = append(taskActionsLog, "successfully call safe delete user auth data")
				return nil
			}
		},
		func() error {
			return nil
		},
	)

	job := NewRetryableRevertibleDeleteUserDataJob(userID, usersClient, authClient, ordersClient, 3, 3)

	q := queue.NewInMemQueue(2)

	jobResultStorage := storage.NewInMemJobResultsStorage()
	taskResultStorage := storage.NewInMemTaskResultStorage()

	p := processor.NewV1Processor(q, jobResultStorage, taskResultStorage)

	err := p.AddJob(job)
	assert.NoError(t, err)

	err = p.Start()
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	err = p.Stop()
	assert.NoError(t, err)

	// check job result
	jobRes, err := p.GetJobResult(job.GetID())
	assert.NoError(t, err)
	assert.Equal(
		t,
		goodjob.JobResult{
			ID:    job.GetID(),
			Value: nil,
			Err:   nil,
		},
		*jobRes,
	)

	// check queue logs
	assert.Equal(
		t,
		[]string{
			"successfully call safe delete user auth data",
			"successfully call safe delete user",
			"failed call safe delete user orders data",
			"failed call safe delete user orders data",
			"failed call safe delete user orders data",
			"failed call safe delete user orders data",
			"failed call revert user data",
			"failed call revert user data",
			"failed call revert user data",
			"successfully call revert user data",
			"failed call revert user auth data",
			"failed call revert user auth data",
			"failed call revert user auth data",
			"successfully call revert user auth data",
		},
		taskActionsLog,
	)
}

func TestRetryableRevertibleDeleteUser_Failure_CantRestoreAuthData(t *testing.T) {
	userID := "user123"

	forwardRetryCountOrders := 4
	backwardRetryCountUsers := 2
	backwardRetryCountAuthData := 4

	var taskActionsLog []string

	usersClient := dummyservices.NewUserServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call safe delete user")
			return nil
		},
		func() error {
			if backwardRetryCountUsers > 0 {
				taskActionsLog = append(taskActionsLog, "failed call revert user data")
				backwardRetryCountUsers--
				return fmt.Errorf("cant revert user data")
			} else {
				taskActionsLog = append(taskActionsLog, "successfully call revert user data")
				return nil
			}
		},
	)

	authClient := dummyservices.NewAuthServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call safe delete user auth data")
			return nil
		},
		func() error {
			if backwardRetryCountAuthData > 0 {
				taskActionsLog = append(taskActionsLog, "failed call revert user auth data")
				backwardRetryCountAuthData--
				return fmt.Errorf("cant revert user auth data")
			} else {
				taskActionsLog = append(taskActionsLog, "successfully call revert user auth data")
				return nil
			}
		},
	)

	ordersClient := dummyservices.NewOrdersServiceClient(
		func() error {
			if forwardRetryCountOrders > 0 {
				taskActionsLog = append(taskActionsLog, "failed call safe delete user orders data")
				forwardRetryCountOrders--
				return fmt.Errorf("cant delete user orders data")
			} else {
				taskActionsLog = append(taskActionsLog, "successfully call safe delete user auth data")
				return nil
			}
		},
		func() error {
			return nil
		},
	)

	job := NewRetryableRevertibleDeleteUserDataJob(userID, usersClient, authClient, ordersClient, 3, 2)

	q := queue.NewInMemQueue(2)

	jobResultStorage := storage.NewInMemJobResultsStorage()
	taskResultStorage := storage.NewInMemTaskResultStorage()

	p := processor.NewV1Processor(q, jobResultStorage, taskResultStorage)

	err := p.AddJob(job)
	assert.NoError(t, err)

	err = p.Start()
	assert.NoError(t, err)

	time.Sleep(3 * time.Second)

	err = p.Stop()
	assert.NoError(t, err)

	// check job result
	jobRes, err := p.GetJobResult(job.GetID())
	assert.NoError(t, err)
	assert.Equal(
		t,
		goodjob.JobResult{
			ID:    job.GetID(),
			Value: nil,
			Err:   fmt.Errorf("cant revert user auth data"),
		},
		*jobRes,
	)

	// check queue logs
	assert.Equal(
		t,
		[]string{
			"successfully call safe delete user auth data",
			"successfully call safe delete user",
			"failed call safe delete user orders data",
			"failed call safe delete user orders data",
			"failed call safe delete user orders data",
			"failed call safe delete user orders data",
			"failed call revert user data",
			"failed call revert user data",
			"successfully call revert user data",
			"failed call revert user auth data",
			"failed call revert user auth data",
			"failed call revert user auth data",
		},
		taskActionsLog,
	)
}
