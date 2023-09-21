package storage

import (
	"fmt"
	"github.com/mtvarkovsky/goodjob/pkg/goodjob"
	"sync"
)

type (
	InMemTaskResultsStorage struct {
		sync.Mutex
		data map[goodjob.JobID]map[goodjob.TaskID]*goodjob.TaskResult
	}

	InMemJobResultsStorage struct {
		sync.Mutex
		data map[goodjob.JobID]*goodjob.JobResult
	}
)

var _ goodjob.TaskResultsStorage = (*InMemTaskResultsStorage)(nil)

func NewInMemTaskResultStorage() goodjob.TaskResultsStorage {
	return &InMemTaskResultsStorage{
		data: make(map[goodjob.JobID]map[goodjob.TaskID]*goodjob.TaskResult),
	}
}

func (s *InMemTaskResultsStorage) Put(result *goodjob.TaskResult) error {
	s.Lock()
	defer s.Unlock()
	if val, ok := s.data[result.JobID]; !ok {
		if val == nil {
			s.data[result.JobID] = make(map[goodjob.TaskID]*goodjob.TaskResult)
		}
	}

	s.data[result.JobID][result.TaskID] = result
	return nil
}

func (s *InMemTaskResultsStorage) Get(jobID goodjob.JobID, taskID goodjob.TaskID) (*goodjob.TaskResult, error) {
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

var _ goodjob.JobResultsStorage = (*InMemJobResultsStorage)(nil)

func NewInMemJobResultsStorage() goodjob.JobResultsStorage {
	return &InMemJobResultsStorage{
		data: make(map[goodjob.JobID]*goodjob.JobResult),
	}
}

func (s *InMemJobResultsStorage) Put(result *goodjob.JobResult) error {
	s.Lock()
	defer s.Unlock()
	s.data[result.ID] = result
	return nil
}

func (s *InMemJobResultsStorage) Get(id goodjob.JobID) (*goodjob.JobResult, error) {
	s.Lock()
	defer s.Unlock()
	if res, found := s.data[id]; found {
		return res, nil
	}

	return nil, fmt.Errorf("job result not found")
}
