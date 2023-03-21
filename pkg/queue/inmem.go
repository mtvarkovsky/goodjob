package queue

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/pkg/interfaces"
	"sync"
)

type (
	// InMemQueue - simple in memory queue without priority. Useful for tests and examples.
	InMemQueue struct {
		mu           sync.Mutex
		maxSize      int
		items        []interfaces.Job
		itemPosition map[interfaces.JobID]int
	}
)

var _ interfaces.Queue = (*InMemQueue)(nil)

func NewInMemQueue(maxSize int) *InMemQueue {
	return &InMemQueue{
		maxSize:      maxSize,
		itemPosition: make(map[interfaces.JobID]int),
	}
}

func (q *InMemQueue) AddJob(job interfaces.Job, args ...any) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == q.maxSize {
		return fmt.Errorf("queue size limit reached")
	}
	if _, found := q.itemPosition[job.GetID()]; found {
		return fmt.Errorf("job with with id=%s already in queue", job.GetID())
	}
	nextPos := len(q.items)
	q.itemPosition[job.GetID()] = nextPos
	q.items = append(q.items, job)
	return nil
}

func (q *InMemQueue) GetNextJob(args ...any) (interfaces.Job, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return nil, nil
	}

	var next interfaces.Job

	for i, job := range q.items {
		if job.GetVisible() {
			next = q.items[i]
			pos, found := q.itemPosition[next.GetID()]
			if !found {
				return nil, fmt.Errorf("job not found in queue")
			}
			q.items[pos].SetVisible(false)
			break
		}
	}

	return next, nil
}

func (q *InMemQueue) RemoveJob(id interfaces.JobID, args ...any) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	pos, found := q.itemPosition[id]
	if !found {
		return fmt.Errorf("job not found in queue")
	}

	if pos < 0 || pos > len(q.items) {
		return fmt.Errorf("position out of bounds")
	}

	for jobID, position := range q.itemPosition {
		if position > pos {
			q.itemPosition[jobID] = position - 1
		}
	}
	delete(q.itemPosition, id)

	if pos == 0 && len(q.items) == 1 {
		q.items = nil
	} else if pos == len(q.items)-1 && len(q.items) == 1 {
		q.items = nil
	} else {
		q.items = append(q.items[:pos], q.items[pos+1:]...)
	}

	return nil
}

func (q *InMemQueue) SetJobVisibility(id interfaces.JobID, visible bool, args ...any) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	pos, found := q.itemPosition[id]
	if !found {
		return fmt.Errorf("job not found in queue")
	}

	q.items[pos].SetVisible(visible)
	return nil
}
