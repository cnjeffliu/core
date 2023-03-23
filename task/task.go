/*
 * @Author: cnzf1
 * @Date: 2023-03-27 19:23:22
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-27 19:23:22
 * @Description:
 */
package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/cnzf1/gocore/funx"
	"github.com/cnzf1/gocore/lang"
	"github.com/cnzf1/gocore/thread"
	"github.com/cnzf1/gocore/timex"
)

type JobStatus int8

//go:generate stringer -type JobStatus -linecomment -output task_string.go
const (
	JOB_STATUS_FAILED  JobStatus = -1 // 失败
	JOB_STATUS_TIMEOUT JobStatus = -2 // 超时
	JOB_STATUS_CREATED JobStatus = 0  // 已创建
	JOB_STATUS_RUNNING JobStatus = 1  // 运行中
	JOB_STATUS_SUCCESS JobStatus = 2  // 成功
)

type JobFunc func(jobID string, in lang.AnyType, out *lang.AnyType) error

type jobGroupBuilder struct {
	group *JobGroup
}

type JobGroup struct {
	GroupID    string
	keyFn      KeyGenerator
	store      Store
	CreateTime string
	Jobs       []*Job
	JobCount   int
	Sync       bool
}

type Job struct {
	JobID     string
	GroupID   string
	BeginTime string
	EndTime   string
	Timeout   time.Duration
	Seq       int
	Status    JobStatus
	mu        sync.Mutex
	fn        JobFunc
	FnIn      lang.AnyType
	FnOut     lang.AnyType
}

func NewJobTaskBuilder(store Store) *jobGroupBuilder {
	if store == nil {
		return nil
	}

	builder := &jobGroupBuilder{
		group: &JobGroup{
			keyFn:    NewTimeBasedRandomGenerator(32),
			store:    store,
			JobCount: 0,
		},
	}

	return builder
}

func (t *jobGroupBuilder) WithKeyFunc(fx KeyGenerator) *jobGroupBuilder {
	t.group.keyFn = fx
	return t
}
func (t *jobGroupBuilder) WithSync() *jobGroupBuilder {
	t.group.Sync = true
	return t
}

func (t *jobGroupBuilder) Build() *JobGroup {
	t.group.GroupID = string(t.group.keyFn.Generate())
	t.group.CreateTime = timex.NowStr(timex.TIME_LAYOUT_COMPACT_MILLSECOND)

	return t.group
}

func (t *JobGroup) CreateJob(fx JobFunc, fnin lang.AnyType, timeout time.Duration) (jobID string) {
	t.JobCount++
	job := &Job{
		JobID:   string(t.keyFn.Generate()),
		GroupID: t.GroupID,
		Timeout: timeout,
		Seq:     t.JobCount,
		Status:  JOB_STATUS_CREATED,
		fn:      fx,
		FnIn:    fnin,
	}
	t.Jobs = append(t.Jobs, job)
	t.store.Add(job.GroupID, job.JobID, job)
	job.UpdateStatus(t.store, JOB_STATUS_CREATED)
	return job.JobID
}

func (t *JobGroup) Run() {
	rg := thread.NewRoutineGroup()
	rg.RunSafe(func() {
		for _, job := range t.Jobs {
			job.Run(t.store)
		}
	})

	if t.Sync {
		rg.Wait()
	}
}

func (t *JobGroup) Get(jobID string) (lang.AnyType, bool) {
	return t.store.Get(jobID)
}

func (t *JobGroup) GetAll(grpID string) []lang.AnyType {
	return t.store.GetAll(grpID)
}

func (t *JobGroup) Keys() []string {
	ids := []string{}
	for _, v := range t.Jobs {
		ids = append(ids, v.JobID)
	}
	return ids
}

func (j *Job) RunLocked(store Store) {
	status := store.GetStatus(j.JobID)
	if status != int(JOB_STATUS_CREATED) {
		fmt.Println("existed")
		return
	}

	store.UpdateStatus(j.JobID, int(JOB_STATUS_RUNNING))
	j.BeginTime = timex.NowStr(timex.TIME_LAYOUT_COMPACT_MILLSECOND)
	err := funx.DoWithTimeout(func() error {
		_err := j.fn(j.JobID, j.FnIn, &j.FnOut)
		if _err != nil {
			fmt.Println(_err)
			return _err
		}
		return nil
	}, funx.WithTimeout(j.Timeout))

	j.EndTime = timex.NowStr(timex.TIME_LAYOUT_COMPACT_MILLSECOND)

	if err == nil {
		j.UpdateStatus(store, JOB_STATUS_SUCCESS)
	} else if err == funx.ErrTimeout {
		j.UpdateStatus(store, JOB_STATUS_TIMEOUT)
		return
	} else {
		j.UpdateStatus(store, JOB_STATUS_FAILED)
	}

	if j.FnOut != nil {
		store.Update(j.JobID, j)
	}
}

func (j *Job) Run(store Store) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.RunLocked(store)
}

func (j *Job) UpdateStatus(store Store, status JobStatus) {
	j.Status = status
	store.UpdateStatus(j.JobID, int(status))
}
