/*
 * @Author: Jeffrey Liu
 * @Date: 2022-08-26 09:26:55
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-11-25 16:56:47
 * @Description:
 */
package timex_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cnzf1/gocore/timex"
)

func TestJitterdBackoff(t *testing.T) {
	stopCh := make(chan struct{})
	defer close(stopCh)
	backoff := timex.NewJitteredBackoffManager(time.Minute, 0.5, timex.RealClock{})

	timex.BackoffUntil(func() {
		fmt.Println("process ", time.Now())
	}, backoff, true, stopCh)
}
