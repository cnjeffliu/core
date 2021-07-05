package limit

import "time"

var slidingWin map[string][]int64
var ok bool

//单机滑动窗口限流
func SlidingWinSingle(key string, count uint, period int64) bool {
	now := time.Now().Unix()
	if slidingWin == nil {
		slidingWin = make(map[string][]int64)
	}

	if _, ok = slidingWin[key]; !ok {
		slidingWin[key] = make([]int64, 0)
	}

	if uint(len(slidingWin[key])) < count {
		slidingWin[key] = append(slidingWin[key], now)
		return true
	}

	firstTime := slidingWin[key][0]

	if now-firstTime <= period {
		return false
	} else {
		slidingWin[key] = slidingWin[key][1:]
		slidingWin[key] = append(slidingWin[key], now)
	}
	return true
}
