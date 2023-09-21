package revertible

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

func TestRevertibleDeleteUser_Success(t *testing.T) {
	userID := "user123"

	var taskActionsLog []string

	usersClient := dummyservices.NewUserServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call safe delete user")
			return nil
		},
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call restore user")
			return nil
		},
	)

	authClient := dummyservices.NewAuthServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call safe user auth data")
			return nil
		},
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call restore user auth data")
			return nil
		},
	)

	ordersClient := dummyservices.NewOrdersServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "failed call safe delete user orders data")
			return fmt.Errorf("cant delete user orders")
		},
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call restore user orders data")
			return nil
		},
	)

	job := NewRevertibleDeleteUserDataJob(userID, usersClient, authClient, ordersClient)

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
			"successfully call safe user auth data",
			"successfully call safe delete user",
			"failed call safe delete user orders data",
			"successfully call restore user",
			"successfully call restore user auth data",
		},
		taskActionsLog,
	)
}

func TestRevertibleDeleteUser_Failure_CantRestoreUser(t *testing.T) {
	userID := "user123"

	var taskActionsLog []string

	usersClient := dummyservices.NewUserServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call safe delete user")
			return nil
		},
		func() error {
			taskActionsLog = append(taskActionsLog, "failed call restore user")
			return fmt.Errorf("cant restore user")
		},
	)

	authClient := dummyservices.NewAuthServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call safe user auth data")
			return nil
		},
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call restore user auth data")
			return nil
		},
	)

	ordersClient := dummyservices.NewOrdersServiceClient(
		func() error {
			taskActionsLog = append(taskActionsLog, "failed call safe delete user orders data")
			return fmt.Errorf("cant delete user orders")
		},
		func() error {
			taskActionsLog = append(taskActionsLog, "successfully call restore user orders data")
			return nil
		},
	)

	job := NewRevertibleDeleteUserDataJob(userID, usersClient, authClient, ordersClient)

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
			Err:   fmt.Errorf("cant restore user"),
		},
		*jobRes,
	)

	// check queue logs
	assert.Equal(
		t,
		[]string{
			"successfully call safe user auth data",
			"successfully call safe delete user",
			"failed call safe delete user orders data",
			"failed call restore user",
		},
		taskActionsLog,
	)
}
