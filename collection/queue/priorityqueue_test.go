package queue_test

import (
	"container/heap"
	"math/rand"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"testing"

	"github.com/cnzf1/gocore/collection/queue"
)

func equal(t *testing.T, act, exp interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n",
			filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}

func TestPriorityQueue(t *testing.T) {
	c := 100
	pq := queue.NewPriorityQueue(c)

	for i := 0; i < c+1; i++ {
		heap.Push(&pq, &queue.PriorityQueueItem{Value: i, Priority: int64(i)})
	}
	equal(t, pq.Len(), c+1)
	equal(t, cap(pq), c*2)

	for i := 0; i < c+1; i++ {
		item := heap.Pop(&pq)
		equal(t, item.(*queue.PriorityQueueItem).Value.(int), i)
	}
	equal(t, cap(pq), c/4)
}

func TestUnsortedInsert(t *testing.T) {
	c := 100
	pq := queue.NewPriorityQueue(c)
	ints := make([]int, 0, c)

	for i := 0; i < c; i++ {
		v := rand.Int()
		ints = append(ints, v)
		heap.Push(&pq, &queue.PriorityQueueItem{Value: i, Priority: int64(v)})
	}
	equal(t, pq.Len(), c)
	equal(t, cap(pq), c)

	sort.Ints(ints)

	for i := 0; i < c; i++ {
		item, _ := pq.PeekAndShift(int64(ints[len(ints)-1]))
		equal(t, item.Priority, int64(ints[i]))
	}
}

func TestRemove(t *testing.T) {
	c := 100
	pq := queue.NewPriorityQueue(c)

	for i := 0; i < c; i++ {
		v := rand.Int()
		heap.Push(&pq, &queue.PriorityQueueItem{Value: "test", Priority: int64(v)})
	}

	for i := 0; i < 10; i++ {
		heap.Remove(&pq, rand.Intn((c-1)-i))
	}

	lastPriority := heap.Pop(&pq).(*queue.PriorityQueueItem).Priority
	for i := 0; i < (c - 10 - 1); i++ {
		item := heap.Pop(&pq)
		equal(t, lastPriority < item.(*queue.PriorityQueueItem).Priority, true)
		lastPriority = item.(*queue.PriorityQueueItem).Priority
	}
}
