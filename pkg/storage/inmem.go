package storage

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/pkg/interfaces"
	"sync"
)

type (
	InMemTaskResultsStorage struct {
		sync.Mutex
		data map[interfaces.JobID]map[interfaces.TaskID]*interfaces.TaskResult
	}

	InMemJobResultsStorage struct {
		sync.Mutex
		data map[interfaces.JobID]*interfaces.JobResult
	}
)

var _ interfaces.TaskResultsStorage = (*InMemTaskResultsStorage)(nil)

func NewInMemTaskResultStorage() interfaces.TaskResultsStorage {
	return &InMemTaskResultsStorage{
		data: make(map[interfaces.JobID]map[interfaces.TaskID]*interfaces.TaskResult),
	}
}

func (s *InMemTaskResultsStorage) Put(result *interfaces.TaskResult) error {
	s.Lock()
	defer s.Unlock()
	if val, ok := s.data[result.JobID]; !ok {
		if val == nil {
			s.data[result.JobID] = make(map[interfaces.TaskID]*interfaces.TaskResult)
		}
	}

	s.data[result.JobID][result.TaskID] = result
	return nil
}

func (s *InMemTaskResultsStorage) Get(jobID interfaces.JobID, taskID interfaces.TaskID) (*interfaces.TaskResult, error) {
	s.Lock()
	defer s.Unlock()
	if _, foundJob := s.data[jobID]; foundJob {
		if res, foundTask := s.data[jobID][taskID]; foundTask {
			return res, nil
		}

		return nil, fmt.Errorf("task result not found")
	}

	return nil, fmt.Errorf("task result for job not found")
}

var _ interfaces.JobResultsStorage = (*InMemJobResultsStorage)(nil)

func NewInMemJobResultsStorage() interfaces.JobResultsStorage {
	return &InMemJobResultsStorage{
		data: make(map[interfaces.JobID]*interfaces.JobResult),
	}
}

func (s *InMemJobResultsStorage) Put(result *interfaces.JobResult) error {
	s.Lock()
	defer s.Unlock()
	s.data[result.ID] = result
	return nil
}

func (s *InMemJobResultsStorage) Get(id interfaces.JobID) (*interfaces.JobResult, error) {
	s.Lock()
	defer s.Unlock()
	if res, found := s.data[id]; found {
		return res, nil
	}

	return nil, fmt.Errorf("job result not found")
}
