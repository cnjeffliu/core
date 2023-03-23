/*
 * @Author: cnzf1
 * @Date: 2023-03-27 19:23:22
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-30 17:34:31
 * @Description:
 */
package task

import (
	"time"

	"github.com/cnzf1/gocore/collection/mapx"
	"github.com/cnzf1/gocore/collection/set"
	"github.com/cnzf1/gocore/lang"
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
	cache  *mapx.ExpiredMap           // key jobid
	group  map[string]*set.Set        // key groupid
	status map[string]*localStoreItem // key jobid
}

type localStoreItem struct {
	grpID  string
	status int
}

func NewLocalStore() Store {
	s := &LocalStore{
		group:  make(map[string]*set.Set),
		status: make(map[string]*localStoreItem),
	}
	fn := func(key string, val lang.AnyType) {
		if v, ok := s.status[key]; ok {
			s.group[v.grpID].Remove(key)
		}
		delete(s.status, key)
	}

	s.cache = mapx.NewExpiredMap(mapx.WithDelCallback(fn))
	return s
}

func (s *LocalStore) Add(grpID string, jobID string, data lang.AnyType) error {
	if _, ok := s.group[grpID]; !ok {
		s.group[grpID] = set.NewSet()
	}
	s.group[grpID].Add(jobID)
	s.cache.Set(jobID, data, 15*time.Minute)
	s.status[jobID] = &localStoreItem{grpID: grpID, status: 0}
	return nil
}

func (s *LocalStore) GetAll(grpID string) []lang.AnyType {
	val, ok := s.group[grpID]
	if !ok {
		return []lang.AnyType{}
	}

	keys := val.Keys()
	var vals []lang.AnyType
	for _, key := range keys {
		v, ok := s.Get(key.(string))
		if ok {
			vals = append(vals, v)
		}
	}
	return vals
}

func (s *LocalStore) Get(jobID string) (lang.AnyType, bool) {
	val, ok := s.cache.Get(jobID)
	return val, ok
}

func (s *LocalStore) GetStatus(jobID string) int {
	if val, ok := s.status[jobID]; ok {
		return val.status
	}
	return 0
}

func (s *LocalStore) Update(jobID string, data lang.AnyType) error {
	s.cache.Set(jobID, data, 15*time.Minute)
	return nil
}

func (s *LocalStore) UpdateStatus(jobID string, status int) error {
	s.status[jobID].status = status
	return nil
}

func (s *LocalStore) Delete(jobID string) error {
	s.cache.Delete(jobID)
	return nil
}

func (s *LocalStore) Clear() error {
	s.cache.Clear()
	for k := range s.status {
		delete(s.status, k)
	}
	for k := range s.group {
		delete(s.group, k)
	}
	return nil
}
