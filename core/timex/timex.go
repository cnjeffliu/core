package timex

import "time"

func NowS() int64 {
	return time.Now().Unix()
}

func NowMS() int64 {
	return time.Now().UnixNano() / 1e6
}

func ElapseMS(begin time.Time) int64 {
	return time.Now().Sub(begin).Microseconds()
}

func ElapseNS(begin time.Time) int64 {
	return time.Now().Sub(begin).Nanoseconds()
}
