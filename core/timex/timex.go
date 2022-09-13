/*
 * @Author: Jeffrey Liu <zhifeng172@163.com>
 * @Date: 2022-07-20 13:56:45
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-09-13 20:41:19
 * @Description:
 */
package timex

import (
	"time"
)

const (
	TIME_LAYOUT_SECOND      = "2006-01-02 15:04:05"
	TIME_LAYOUT_MILLSECOND  = "2006-01-02 15:04:05.000"
	TIME_LAYOUT_MICROSECOND = "2006-01-02 15:04:05.000000"
	TIME_LAYOUT_NANOSECOND  = "2006-01-02 15:04:05.000000000"
)

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

// input format is 2022-01-01 01:00:00
func StrToTime(s string) time.Time {
	d, err := time.ParseInLocation(TIME_LAYOUT_SECOND, s, time.Local)
	if err != nil {
		return time.Now()
	}

	return d
}

// output format is 2022-01-01 01:00:00
func TimeToStr(s time.Time) string {
	return s.Format(TIME_LAYOUT_SECOND)
}

func TSToStr(sec int64, nsec int64) string {
	return time.Unix(sec, nsec).Format(TIME_LAYOUT_SECOND)
}
