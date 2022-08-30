/*
 * @Author: Jeffrey Liu <zhifeng172@163.com>
 * @Date: 2021-12-15 14:18:20
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-08-27 16:06:04
 * @Description:
 */
package thread

import "sync"

type RoutineGroup struct {
	waitGroup sync.WaitGroup
}

func NewRoutineGroup() *RoutineGroup {
	return new(RoutineGroup)
}

func (g *RoutineGroup) Run(fn func()) {
	g.waitGroup.Add(1)

	go func() {
		defer g.waitGroup.Done()
		fn()
	}()
}

func (g *RoutineGroup) RunSafe(fn func()) {
	g.waitGroup.Add(1)

	GoSafe(func() {
		defer g.waitGroup.Done()
		fn()
	})
}

func (g *RoutineGroup) Wait() {
	g.waitGroup.Wait()
}
