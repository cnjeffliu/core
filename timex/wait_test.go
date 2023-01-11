/*
 * @Author: cnzf1
 * @Date: 2022-08-26 09:26:55
 * @LastEditors: cnzf1
 * @LastEditTime: 2022-12-15 16:13:53
 * @Description:
 */
package timex_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cnzf1/gocore/timex"
	"github.com/stretchr/testify/assert"
)

func TestJitterUp(t *testing.T) {
	s := time.Hour
	d := timex.JitterUp(s, 1)

	assert.GreaterOrEqual(t, d.Seconds(), s.Seconds())
	assert.LessOrEqual(t, d.Seconds(), 2*s.Seconds())
}

func TestJitterAround(t *testing.T) {
	s := time.Hour
	d := timex.JitterAround(s, 1)

	assert.GreaterOrEqual(t, d.Seconds(), 0.0)
	assert.LessOrEqual(t, d.Seconds(), 2*s.Seconds())
}
func TestJitterdBackoff(t *testing.T) {
	stopCh := make(chan struct{})
	defer close(stopCh)
	backoff := timex.NewJitteredBackoffManager(time.Minute, 0.5, timex.RealClock{})

	timex.BackoffUntil(func() {
		fmt.Println("process ", time.Now())
	}, backoff, true, stopCh)
}
