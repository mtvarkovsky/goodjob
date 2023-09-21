package retryable

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

func TestRetryableDeleteUser_Success(t *testing.T) {
	userID := "user123"

	var taskActionsLog []string

	retryCount := 3

	usersClient := dummyservices.NewUserServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call safe delete user")
			return nil
		},
		func() error {
			return nil
		},
	)

	authClient := dummyservices.NewAuthServiceClient(
		func() error {
			if retryCount > 0 {
				taskActionsLog = append(taskActionsLog, "failed call safe delete user auth data")
				retryCount--
				return fmt.Errorf("cant delete user auth data")
			} else {
				taskActionsLog = append(taskActionsLog, "successfully call safe delete user auth data")
				return nil
			}
		},
		func() error {
			return nil
		},
	)

	ordersClient := dummyservices.NewOrdersServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call safe delete user orders data")
			return nil
		},
		func() error {
			return nil
		},
	)

	job := NewRetryableDeleteUserDataJob(userID, usersClient, authClient, ordersClient, 4)

	q := queue.NewInMemQueue(2)

	jobResultStorage := storage.NewInMemJobResultsStorage()
	taskResultStorage := storage.NewInMemTaskResultStorage()

	p := processor.NewV1Processor(q, jobResultStorage, taskResultStorage)

	err := p.AddJob(job)
	assert.NoError(t, err)

	err = p.Start()
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

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
			"failed call safe delete user auth data",
			"failed call safe delete user auth data",
			"failed call safe delete user auth data",
			"successfully call safe delete user auth data",
			"successfully call safe delete user",
			"successfully call safe delete user orders data",
		},
		taskActionsLog,
	)
}

func TestRetryableDeleteUser_Failure_RetryCountExceeded(t *testing.T) {
	userID := "user123"

	var taskActionsLog []string

	usersClient := dummyservices.NewUserServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call safe delete user")
			return nil
		},
		func() error {
			return nil
		},
	)

	authClient := dummyservices.NewAuthServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "failed call safe delete user auth data")
			return fmt.Errorf("cant delete user auth data")
		},
		func() error {
			return nil
		},
	)

	ordersClient := dummyservices.NewOrdersServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call safe delete user orders data")
			return nil
		},
		func() error {
			return nil
		},
	)

	job := NewRetryableDeleteUserDataJob(userID, usersClient, authClient, ordersClient, 4)

	q := queue.NewInMemQueue(2)

	jobResultStorage := storage.NewInMemJobResultsStorage()
	taskResultStorage := storage.NewInMemTaskResultStorage()

	p := processor.NewV1Processor(q, jobResultStorage, taskResultStorage)

	err := p.AddJob(job)
	assert.NoError(t, err)

	err = p.Start()
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

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
			Err:   fmt.Errorf("cant delete user auth data"),
		},
		*jobRes,
	)

	// check queue logs
	assert.Equal(
		t,
		[]string{
			"failed call safe delete user auth data",
			"failed call safe delete user auth data",
			"failed call safe delete user auth data",
			"failed call safe delete user auth data",
			"failed call safe delete user auth data",
		},
		taskActionsLog,
	)
}
