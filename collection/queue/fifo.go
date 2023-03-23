/*
 * @Author: cnzf1
 * @Date: 2022-07-20 13:56:02
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-30 17:11:20
 * @Description:
 */

package queue

import (
	"sync"
	"sync/atomic"

	"github.com/cnzf1/gocore/collection/set"
	"github.com/cnzf1/gocore/lang"
)

func NewFIFOQueue() *FIFOQueue {
	t := &FIFOQueue{
		dirty:      set.NewSet(),
		processing: set.NewSet(),
		cond:       sync.NewCond(&sync.Mutex{}),
	}

	return t
}

type FIFOQueue struct {
	queue      []lang.AnyType
	queueLen   int32
	dirty      *set.Set
	processing *set.Set
	processLen int32

	cond *sync.Cond

	shuttingDown bool
	drain        bool
}

func (q *FIFOQueue) addQueue(delta int) {
	atomic.AddInt32(&q.queueLen, int32(delta))
}

func (q *FIFOQueue) addProcess(delta int) {
	atomic.AddInt32(&q.processLen, int32(delta))
}

func (q *FIFOQueue) Add(item interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if q.shuttingDown {
		return
	}
	if q.dirty.Contains(item) {
		return
	}

	q.dirty.Add(item)
	if q.processing.Contains(item) {
		return
	}

	q.queue = append(q.queue, item)
	q.addQueue(1)
	q.cond.Signal()
}

func (q *FIFOQueue) Len() int {
	return int(atomic.LoadInt32(&q.queueLen))
}

func (q *FIFOQueue) ProcessingLen() int {
	return int(atomic.LoadInt32(&q.processLen))
}

func (q *FIFOQueue) Get() (item interface{}, shutdown bool) {
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

	q.processing.Add(item)
	q.addProcess(1)
	q.dirty.Remove(item)

	return item, false
}

func (q *FIFOQueue) Done(item interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.processing.Remove(item)
	q.addProcess(-1)
	if q.dirty.Contains(item) {
		q.queue = append(q.queue, item)
		q.addQueue(1)
		q.cond.Signal()
	} else if q.processing.Count() == 0 {
		q.cond.Signal()
	}
}

func (q *FIFOQueue) ShutDown() {
	q.setDrain(false)
	q.shutdown()
}

func (q *FIFOQueue) ShutDownWithDrain() {
	q.setDrain(true)
	q.shutdown()
	for q.isProcessing() && q.shouldDrain() {
		q.waitForProcessing()
	}
}

func (q *FIFOQueue) isProcessing() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.processing.Count() != 0
}

func (q *FIFOQueue) waitForProcessing() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if q.processing.Count() == 0 {
		return
	}
	q.cond.Wait()
}

func (q *FIFOQueue) setDrain(shouldDrain bool) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.drain = shouldDrain
}

func (q *FIFOQueue) shouldDrain() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.drain
}

func (q *FIFOQueue) shutdown() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.shuttingDown = true
	q.cond.Broadcast()
}

func (q *FIFOQueue) ShuttingDown() bool {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	return q.shuttingDown
}
