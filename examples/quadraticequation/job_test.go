package quadraticequation

import (
	"github.com/mtvarkovsky/goodjob/pkg/processor"
	"github.com/mtvarkovsky/goodjob/pkg/queue"
	"github.com/mtvarkovsky/goodjob/pkg/storage"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TODO: rework tests to remove time.Sleep()

func TestQuadraticEquation(t *testing.T) {
	job1 := NewQuadraticEquationJob(1.0, 5.0, 6.0)
	job2 := NewQuadraticEquationJob(1.0, -7.0, -3.0)

	q := queue.NewInMemQueue(2)

	jobResultStorage := storage.NewInMemJobResultsStorage()
	taskResultStorage := storage.NewInMemTaskResultStorage()

	p := processor.NewV1Processor(q, jobResultStorage, taskResultStorage)

	err := p.AddJob(job1)
	assert.NoError(t, err)
	err = p.Start()
	assert.NoError(t, err)
	err = p.AddJob(job2)
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	err = p.Stop()
	assert.NoError(t, err)

	job1Res, err := p.GetJobResult(job1.GetID())
	assert.NoError(t, err)
	job2Res, err := p.GetJobResult(job2.GetID())
	assert.NoError(t, err)

	assert.Equal(t, []float64{-2.0, -3.0}, job1Res.Value)
	assert.Equal(t, []float64{7.405124837953327, -0.405124837953327}, job2Res.Value)
}
