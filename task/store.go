package task

import (
	"github.com/cnzf1/gocore/lang"
	"github.com/cnzf1/gocore/setx"
)

type Store interface {
	Add(grpID string, jobID string, data lang.AnyType) error
	Get(jobID string) (lang.AnyType, bool)
	GetAll(grpID string) []lang.AnyType
	GetStatus(jobID string) int
	Update(jobID string, data lang.AnyType) error
	UpdateStatus(jobID string, status int) error
	Delete(jobID string) error
	Clear() error
}

type LocalStore struct {
	cache  map[string]lang.AnyType // key jobid
	group  map[string]setx.String  // key groupid
	status map[string]int          // key jobid
}

func NewLocalStore() Store {
	return &LocalStore{
		cache:  make(map[string]lang.AnyType),
		group:  make(map[string]setx.String),
		status: make(map[string]int),
	}
}

func (s *LocalStore) Add(grpID string, jobID string, data lang.AnyType) error {
	if _, ok := s.group[grpID]; !ok {
		s.group[grpID] = setx.NewString()
	}
	s.group[grpID].Insert(jobID)
	s.cache[jobID] = data
	s.status[jobID] = 0
	return nil
}

func (s *LocalStore) GetAll(grpID string) []lang.AnyType {
	val, ok := s.group[grpID]
	if !ok {
		return []lang.AnyType{}
	}

	keys := val.List()
	var vals []lang.AnyType
	for _, key := range keys {
		v, ok := s.Get(key)
		if ok {
			vals = append(vals, v)
		}
	}
	return vals
}

func (s *LocalStore) Get(jobID string) (lang.AnyType, bool) {
	val, ok := s.cache[jobID]
	return val, ok
}

func (s *LocalStore) GetStatus(jobID string) int {
	if val, ok := s.status[jobID]; ok {
		return val
	}
	return 0
}

func (s *LocalStore) Update(jobID string, data lang.AnyType) error {
	s.cache[jobID] = data
	return nil
}

func (s *LocalStore) UpdateStatus(jobID string, status int) error {
	s.status[jobID] = status
	return nil
}

func (s *LocalStore) Delete(jobID string) error {
	delete(s.cache, jobID)
	delete(s.status, jobID)
	return nil
}

func (s *LocalStore) Clear() error {
	for k := range s.cache {
		delete(s.cache, k)
	}
	for k := range s.status {
		delete(s.status, k)
	}
	return nil
}
