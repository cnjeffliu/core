/*
 * @Author: Jeffrey Liu
 * @Date: 2022-12-12 16:21:15
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-12-12 16:38:31
 * @Description:
 */
package queue_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/cnzf1/gocore/lang"
	"github.com/cnzf1/gocore/queue"
)

func TestEventHandler(t *testing.T) {
	f := queue.HandlerFunc(func(s lang.AnyType) {
		fmt.Println("handle  ", lang.Repr(s))
	})

	e := queue.NewEventHandler(f, 100)

	for i := 0; i < 100; i++ {
		fmt.Println("request ", strconv.Itoa(i))
		e.Push(i)
	}

	time.Sleep(time.Second)
	e.Stop()
	e.Stop()
}
