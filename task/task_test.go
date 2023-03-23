/*
 * @Author: cnzf1
 * @Date: 2023-03-27 19:23:22
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-27 19:23:22
 * @Description:
 */
package task_test

import (
	"testing"
	"time"

	"github.com/cnzf1/gocore/lang"
	"github.com/cnzf1/gocore/task"
	"github.com/stretchr/testify/assert"
)

func TestJobTask_RunSync(t *testing.T) {
	store := task.NewLocalStore()
	jobTask := task.NewJobTaskBuilder(store).WithSync().Build()
	jobid := jobTask.CreateJob(func(jobID string, in lang.AnyType, out *lang.AnyType) error {
		time.Sleep(3 * time.Second)
		*out = "this is output"
		return nil
	}, nil, 5*time.Second)
	jobTask.Run()

	assert.Equal(t, task.JOB_STATUS_SUCCESS, task.JobStatus(store.GetStatus(jobid)))
	job, _ := store.Get(jobid)
	assert.Equal(t, "this is output", job.(*task.Job).FnOut)
}

func TestJobTask_Run(t *testing.T) {
	store := task.NewLocalStore()
	jobTask := task.NewJobTaskBuilder(store).Build()
	jobid := jobTask.CreateJob(func(jobID string, in lang.AnyType, out *lang.AnyType) error {
		time.Sleep(3 * time.Second)
		*out = "this is output"
		return nil
	}, nil, 5*time.Second)
	jobTask.Run()

	time.Sleep(5 * time.Second)
	assert.Equal(t, task.JOB_STATUS_SUCCESS, task.JobStatus(store.GetStatus(jobid)))
	job, _ := store.Get(jobid)
	assert.Equal(t, "this is output", job.(*task.Job).FnOut)
}

func TestJobTask_RunTimeout(t *testing.T) {
	store := task.NewLocalStore()
	jobTask := task.NewJobTaskBuilder(store).Build()
	jobid := jobTask.CreateJob(func(jobID string, in lang.AnyType, out *lang.AnyType) error {
		time.Sleep(5 * time.Second)
		*out = "this is output"
		return nil
	}, nil, 3*time.Second)
	jobTask.Run()

	time.Sleep(5 * time.Second)
	assert.Equal(t, task.JOB_STATUS_TIMEOUT, task.JobStatus(store.GetStatus(jobid)))
}
