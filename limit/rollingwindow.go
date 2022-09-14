/*
* 滑动窗口限流
* Jeff.Liu 2021.07.12
**/
package limit

import (
	"sync"

	"github.com/cnjeffliu/gocore/timex"
)

type WindowLimit struct {
	lock   sync.RWMutex
	limit  uint64
	period int64
	win    []int64
}

const (
	defaultLimit  = 500
	defaultPeriod = 60 * 1e3 // 60s
)

type WindowLimitOption func(w *WindowLimit)

func WithLimit(cnt uint64) WindowLimitOption {
	return func(w *WindowLimit) {
		w.limit = cnt
	}
}

// period sec
func WithPeriod(period int64) WindowLimitOption {
	return func(w *WindowLimit) {
		w.period = period * 1e3
	}
}

/*
max: limit num in a period
period: second
*/
func NewWindowLimit(opts ...WindowLimitOption) *WindowLimit {
	wl := &WindowLimit{
		limit:  defaultLimit,
		period: defaultPeriod,
		win:    make([]int64, 0),
	}

	for _, opt := range opts {
		opt(wl)
	}

	return wl
}

func (wl *WindowLimit) Access() bool {
	now := timex.NowMS()

	wl.lock.Lock()
	defer wl.lock.Unlock()

	if uint64(len(wl.win)) < wl.limit {
		wl.win = append(wl.win, now)
		return true
	}

	if now-wl.win[0] < wl.period {
		return false
	}

	wl.win = wl.win[1:]
	wl.win = append(wl.win, now)

	return true
}

func (wl *WindowLimit) Count() uint64 {
	now := timex.NowMS()
	wl.lock.RLock()
	defer wl.lock.RUnlock()

	idx := uint64(0)
	for k, v := range wl.win {
		if now-v < wl.period {
			idx = uint64(k)
			break
		}
	}

	return uint64(len(wl.win)) - idx
}
