/*
 * @Author: cnzf1
 * @Date: 2023-01-04 10:16:05
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-30 17:15:18
 * @Description:
 */
package queue

import (
	"container/heap"

	"github.com/cnzf1/gocore/lang"
)

type PriorityQueueItem struct {
	Value    lang.AnyType
	Priority int64
	Index    int
}

// this is a priority queue as implemented by a min heap
// ie. the 0th element is the *lowest* value
type PriorityQueue []*PriorityQueueItem

func NewPriorityQueue(capacity int) PriorityQueue {
	return make(PriorityQueue, 0, capacity)
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	c := cap(*pq)
	if n+1 > c {
		npq := make(PriorityQueue, n, c*2)
		copy(npq, *pq)
		*pq = npq
	}
	*pq = (*pq)[0 : n+1]
	item := x.(*PriorityQueueItem)
	item.Index = n
	(*pq)[n] = item
}

func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	c := cap(*pq)
	if n < (c/2) && c > 25 {
		npq := make(PriorityQueue, n, c/2)
		copy(npq, *pq)
		*pq = npq
	}
	item := (*pq)[n-1]
	item.Index = -1
	*pq = (*pq)[0 : n-1]
	return item
}

func (pq *PriorityQueue) PeekAndShift(max int64) (*PriorityQueueItem, int64) {
	if pq.Len() == 0 {
		return nil, 0
	}

	item := (*pq)[0]
	if item.Priority > max {
		return nil, item.Priority - max
	}
	heap.Remove(pq, 0)

	return item, 0
}
