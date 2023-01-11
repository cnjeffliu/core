/*
 * @Author: cnzf1
 * @Date: 2022-07-20 13:56:02
 * @LastEditors: cnzf1
 * @LastEditTime: 2022-10-21 23:34:09
 * @Description:
 */

package queue

import (
	"sync"
	"sync/atomic"

	"github.com/cnzf1/gocore/lang"
	"github.com/cnzf1/gocore/setx"
)

func NewQueue() *Queue {
	t := &Queue{
		dirty:      setx.Set{},
		processing: setx.Set{},
		cond:       sync.NewCond(&sync.Mutex{}),
	}

	return t
}

type Queue struct {
	queue      []lang.AnyType
	queueLen   int32
	dirty      setx.Set
	processing setx.Set
	processLen int32

	cond *sync.Cond

	shuttingDown bool
	drain        bool
}

func (q *Queue) addQueue(delta int) {
	atomic.AddInt32(&q.queueLen, int32(delta))
}

func (q *Queue) addProcess(delta int) {
	atomic.AddInt32(&q.processLen, int32(delta))
}

func (q *Queue) Add(item interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if q.shuttingDown {
		return
	}
	if q.dirty.Has(item) {
		return
	}

	q.dirty.Insert(item)
	if q.processing.Has(item) {
		return
	}

	q.queue = append(q.queue, item)
	q.addQueue(1)
	q.cond.Signal()
}

func (q *Queue) Len() int {
	return int(atomic.LoadInt32(&q.queueLen))
}

func (q *Queue) ProcessingLen() int {
	return int(atomic.LoadInt32(&q.processLen))
}

func (q *Queue) Get() (item interface{}, shutdown bool) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for len(q.queue) == 0 && !q.shuttingDown {
		q.cond.Wait()
	}
	if len(q.queue) == 0 {
		return nil, true
	}

	item = q.queue[0]
	q.queue[0] = nil
	q.queue = q.queue[1:]
	q.addQueue(-1)

	q.processing.Insert(item)
	q.addProcess(1)
	q.dirty.Delete(item)

	return item, false
}

func (q *Queue) Done(item interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.processing.Delete(item)
	q.addProcess(-1)
	if q.dirty.Has(item) {
		q.queue = append(q.queue, item)
		q.addQueue(1)
		q.cond.Signal()
	} else if q.processing.Len() == 0 {
		q.cond.Signal()
	}
}

func (q *Queue) ShutDown() {
	q.setDrain(false)
	q.shutdown()
}

func (q *Queue) ShutDownWithDrain() {
	q.setDrain(true)
	q.shutdown()
	for q.isProcessing() && q.shouldDrain() {
		q.waitForProcessing()
	}
}

func (q *Queue) isProcessing() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.processing.Len() != 0
}

func (q *Queue) waitForProcessing() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if q.processing.Len() == 0 {
		return
	}
	q.cond.Wait()
}

func (q *Queue) setDrain(shouldDrain bool) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.drain = shouldDrain
}

func (q *Queue) shouldDrain() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.drain
}

func (q *Queue) shutdown() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.shuttingDown = true
	q.cond.Broadcast()
}

func (q *Queue) ShuttingDown() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	return q.shuttingDown
}
