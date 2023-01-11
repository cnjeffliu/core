/*
 * @Author: cnzf1
 * @Date: 2022-12-12 15:22:54
 * @LastEditors: cnzf1
 * @LastEditTime: 2022-12-13 14:45:01
 * @Description:
 */
package queue

import (
	"fmt"
	"sync/atomic"

	"github.com/cnzf1/gocore/lang"
	"github.com/cnzf1/gocore/thread"
)

const defaultSize = 10

type HandlerFunc func(lang.AnyType)

type EventHandler struct {
	quit    chan lang.PlaceholderType
	queue   chan lang.AnyType
	h       HandlerFunc
	started uint32
}

func NewEventHandler(h HandlerFunc, size ...int) *EventHandler {
	n := defaultSize
	if len(size) > 0 && size[0] >= 0 {
		n = size[0]
	}
	q := &EventHandler{}
	q.quit = make(chan lang.PlaceholderType)
	q.queue = make(chan lang.AnyType, n)
	q.h = h

	q.serve()
	return q
}

// Push a request msg to handler
func (e *EventHandler) Push(req lang.AnyType) {
	e.queue <- req
}

func (e *EventHandler) serve() {
	if !atomic.CompareAndSwapUint32(&e.started, 0, 1) {
		fmt.Println("has already called start")
		return
	}

	f := func() {
		for {
			select {
			case req := <-e.queue:
				e.h(req)
			case <-e.quit:
				break
			}
		}
	}

	thread.GoSafe(f)
}

func (e *EventHandler) close() {
	if !atomic.CompareAndSwapUint32(&e.started, 1, 0) {
		// has already called close
		fmt.Println("has already called close")
		return
	}
	e.quit <- lang.Placeholder
}

func (e *EventHandler) Stop() {
	e.close()
}

func (e *EventHandler) Len() int {
	return len(e.queue)
}
