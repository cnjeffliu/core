/*
 * @Author: Jeffrey Liu
 * @Date: 2022-08-26 09:26:55
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-08-26 09:38:55
 * @Description:
 */
package timex_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cnjeffliu/gocore/clock"
	"github.com/cnjeffliu/gocore/timex"
)

func TestJitterdBackoff(t *testing.T) {
	stopCh := make(chan struct{})
	defer close(stopCh)
	backoff := timex.NewJitteredBackoffManager(time.Minute, 0.5, clock.RealClock{})

	timex.BackoffUntil(func() {
		fmt.Println("process ", time.Now())
	}, backoff, true, stopCh)
}
