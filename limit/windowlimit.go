package limit

import (
	"serv/timex"
	"sync"
)

type WindowLimit struct {
	lock   sync.RWMutex
	max    uint64
	period int64
	win    []int64
}

/*
 count: max num in a period
 period: second
*/
func NewWindowLimit(count uint64, period int64) *WindowLimit {
	wl := &WindowLimit{
		max:    count,
		period: period * 1e3,
		win:    make([]int64, 0),
	}

	return wl
}

func (wl *WindowLimit) Access() bool {
	now := timex.NowMS()

	wl.lock.Lock()
	defer wl.lock.Unlock()

	if uint64(len(wl.win)) < wl.max {
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

	idx := wl.max
	for k, v := range wl.win {
		if now-v < wl.period {
			idx = uint64(k)
			break
		}
	}
	return uint64(len(wl.win)) - idx
}
