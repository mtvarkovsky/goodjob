package queue

import (
	"github.com/golang/mock/gomock"
	"github.com/mtvarkovsky/goodjob/mocks"
	"github.com/mtvarkovsky/goodjob/pkg/interfaces"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InMemQueue_AddJob(t *testing.T) {
	ctrl := gomock.NewController(t)

	maxQueueSize := 3
	queue := NewInMemQueueWithLogs(maxQueueSize)

	jobID1 := interfaces.JobID("job1")
	jobID2 := interfaces.JobID("job2")
	jobID3 := interfaces.JobID("job3")

	// try to add first job
	job1 := mocks.NewMockJob(ctrl)
	job1.EXPECT().GetID().Return(jobID1).Times(2)
	err := queue.AddJob(job1)
	assert.NoError(t, err)

	// try to add first job again
	job1.EXPECT().GetID().Return(jobID1).Times(2)
	err = queue.AddJob(job1)
	assert.Error(t, err)

	// try to add second job
	job2 := mocks.NewMockJob(ctrl)
	job2.EXPECT().GetID().Return(jobID2).Times(2)
	err = queue.AddJob(job2)
	assert.NoError(t, err)

	// try to add third job
	job3 := mocks.NewMockJob(ctrl)
	job3.EXPECT().GetID().Return(jobID3).Times(2)
	err = queue.AddJob(job3)
	assert.NoError(t, err)

	// try to add fourth job
	job4 := mocks.NewMockJob(ctrl)
	err = queue.AddJob(job4)
	assert.Error(t, err)
}

func Test_InMemQueue_RemoveJob(t *testing.T) {
	ctrl := gomock.NewController(t)

	maxQueueSize := 3
	queue := NewInMemQueueWithLogs(maxQueueSize)

	jobID1 := interfaces.JobID("job1")
	jobID2 := interfaces.JobID("job2")
	jobID3 := interfaces.JobID("job3")

	// try to add first job
	job1 := mocks.NewMockJob(ctrl)
	job1.EXPECT().GetID().Return(jobID1).Times(2)
	err := queue.AddJob(job1)
	assert.NoError(t, err)

	// try to add second job
	job2 := mocks.NewMockJob(ctrl)
	job2.EXPECT().GetID().Return(jobID2).Times(2)
	err = queue.AddJob(job2)
	assert.NoError(t, err)

	// try to add third job
	job3 := mocks.NewMockJob(ctrl)
	job3.EXPECT().GetID().Return(jobID3).Times(2)
	err = queue.AddJob(job3)
	assert.NoError(t, err)

	// try to remove first job
	err = queue.RemoveJob(jobID1)
	assert.NoError(t, err)

	// try to remove second job
	err = queue.RemoveJob(jobID2)
	assert.NoError(t, err)

	// try to remove third job
	err = queue.RemoveJob(jobID3)
	assert.NoError(t, err)

	// try to remove third job again
	err = queue.RemoveJob(jobID3)
	assert.Error(t, err)
}

func Test_InMemQueue_SetJobVisibility(t *testing.T) {
	ctrl := gomock.NewController(t)

	maxQueueSize := 3
	queue := NewInMemQueueWithLogs(maxQueueSize)

	jobID := interfaces.JobID("job")

	// try to add job
	job := mocks.NewMockJob(ctrl)
	job.EXPECT().GetID().Return(jobID).Times(2)
	err := queue.AddJob(job)
	assert.NoError(t, err)

	// try to set job to invisible
	job.EXPECT().SetVisible(false)
	err = queue.SetJobVisibility(jobID, false)
	assert.NoError(t, err)

	// try to set job to visible
	job.EXPECT().SetVisible(true)
	err = queue.SetJobVisibility(jobID, true)
	assert.NoError(t, err)

	// try to set job visibility for task that does not exist
	err = queue.RemoveJob(jobID)
	assert.NoError(t, err)
	err = queue.SetJobVisibility(jobID, false)
	assert.Error(t, err)
}

func Test_InMemQueue_GetNextJob(t *testing.T) {
	ctrl := gomock.NewController(t)

	maxQueueSize := 3
	queue := NewInMemQueueWithLogs(maxQueueSize)

	// try to get job from empty queue
	j, err := queue.GetNextJob()
	assert.NoError(t, err)
	assert.Nil(t, j)

	jobID1 := interfaces.JobID("job1")
	jobID2 := interfaces.JobID("job2")
	jobID3 := interfaces.JobID("job3")

	// try to add first job
	job1 := mocks.NewMockJob(ctrl)
	job1.EXPECT().GetID().Return(jobID1).Times(2)
	err = queue.AddJob(job1)
	assert.NoError(t, err)

	// try to add second job
	job2 := mocks.NewMockJob(ctrl)
	job2.EXPECT().GetID().Return(jobID2).Times(2)
	err = queue.AddJob(job2)
	assert.NoError(t, err)

	// try to add third job
	job3 := mocks.NewMockJob(ctrl)
	job3.EXPECT().GetID().Return(jobID3).Times(2)
	err = queue.AddJob(job3)
	assert.NoError(t, err)

	// get first job
	job1.EXPECT().GetVisible().Return(true)
	job1.EXPECT().SetVisible(false)
	job1.EXPECT().GetID().Return(jobID1)
	j, err = queue.GetNextJob()
	assert.NoError(t, err)
	assert.NotNil(t, j)
	assert.Equal(t, job1, j)

	// get second job
	job1.EXPECT().GetVisible().Return(false)
	job2.EXPECT().GetVisible().Return(true)
	job2.EXPECT().SetVisible(false)
	job2.EXPECT().GetID().Return(jobID2)
	j, err = queue.GetNextJob()
	assert.NoError(t, err)
	assert.NotNil(t, j)
	assert.Equal(t, job2, j)

	// get third job
	job1.EXPECT().GetVisible().Return(false)
	job2.EXPECT().GetVisible().Return(false)
	job3.EXPECT().GetVisible().Return(true)
	job3.EXPECT().SetVisible(false)
	job3.EXPECT().GetID().Return(jobID3)
	j, err = queue.GetNextJob()
	assert.NoError(t, err)
	assert.NotNil(t, j)
	assert.Equal(t, job3, j)

	// try to get another job
	job1.EXPECT().GetVisible().Return(false)
	job2.EXPECT().GetVisible().Return(false)
	job3.EXPECT().GetVisible().Return(false)
	j, err = queue.GetNextJob()
	assert.NoError(t, err)
	assert.Nil(t, j)
}
