package storage

import (
	"github.com/mtvarkovsky/goodjob/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InMemTaskResultsStorage(t *testing.T) {
	jobID := interfaces.JobID("testJob")
	taskID := interfaces.TaskID("testTask")
	taskResult := &interfaces.TaskResult{
		JobID:  jobID,
		TaskID: taskID,
		Value:  nil,
		Err:    nil,
	}

	storage := NewInMemTaskResultStorage()

	err := storage.Put(taskResult)
	assert.NoError(t, err)

	res, err := storage.Get(jobID, taskID)
	assert.NoError(t, err)
	assert.Equal(t, res, taskResult)

	res, err = storage.Get(jobID, "-")
	assert.Error(t, err)
	assert.Nil(t, res)

	res, err = storage.Get("-", taskID)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func Test_InMemJobResultsStorage(t *testing.T) {
	jobID := interfaces.JobID("testJob")
	jobResult := &interfaces.JobResult{
		ID:    jobID,
		Value: nil,
		Err:   nil,
	}

	storage := NewInMemJobResultsStorage()

	err := storage.Put(jobResult)
	assert.NoError(t, err)

	res, err := storage.Get("-")
	assert.Error(t, err)
	assert.Nil(t, res)
}
