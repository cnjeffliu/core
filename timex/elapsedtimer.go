/*
 * @Author: cnzf1
 * @Date: 2023-01-04 09:58:25
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-01-04 09:58:29
 * @Description:
 */
package timex

import "time"

type ElapsedTimer struct {
	start time.Time
}

// NewElapsedTimer returns an ElapsedTimer.
func NewElapsedTimer() *ElapsedTimer {
	return &ElapsedTimer{
		start: time.Now(),
	}
}

// Duration returns the elapsed time.
func (et *ElapsedTimer) Duration() time.Duration {
	return time.Since(et.start)
}

// Elapsed returns the string representation of elapsed time.
func (et *ElapsedTimer) Elapsed() time.Duration {
	return time.Since(et.start)
}

// ElapsedMs returns the elapsed time of string on milliseconds.
func (et *ElapsedTimer) ElapsedMs() float32 {
	return float32(time.Since(et.start)) / float32(time.Millisecond)
}
