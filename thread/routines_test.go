/*
 * @Author: cnzf1
 * @Date: 2023-01-11 13:57:24
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-01-11 14:00:05
 * @Description:
 */
package thread_test

import (
	"fmt"
	"testing"

	"github.com/cnzf1/gocore/thread"
)

func f1() {
	fmt.Println("hello f1")
}

func TestSafe(t *testing.T) {
	thread.GoSafe(func() {
		fmt.Println("ok")
	})

	f := func() {
		fmt.Println("hello f")
	}

	thread.GoSafe(f)

	thread.GoSafe(f1)
}
