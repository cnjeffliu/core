/*
 * @Author: Jeffrey Liu <zhifeng172@163.com>
 * @Date: 2021-12-15 14:18:20
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-08-30 11:34:24
 * @Description:
 */
package thread

type WorkerGroup struct {
	job     func()
	workers int
}

func NewWorkerGroup(job func(), workers int) WorkerGroup {
	return WorkerGroup{
		job:     job,
		workers: workers,
	}
}

func (wg WorkerGroup) Start() {
	group := NewRoutineGroup()
	for i := 0; i < wg.workers; i++ {
		group.RunSafe(wg.job)
	}
	group.Wait()
}
