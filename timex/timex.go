/*
 * @Author: Jeffrey Liu
 * @Date: 2022-07-20 13:56:45
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-10-24 11:04:57
 * @Description:
 */
package timex

import (
	"strings"
	"time"
)

const (
	TIME_LAYOUT_YEAR        = "2006"
	TIME_LAYOUT_MONTH       = "2006-01"
	TIME_LAYOUT_DAY         = "2006-01-02"
	TIME_LAYOUT_MINUTE      = "2006-01-02 15:04"
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

// output format is 2022-01-01 01:00:00
func TSToStr(sec int64, nsec int64) string {
	return time.Unix(sec, nsec).Format(TIME_LAYOUT_SECOND)
}

// output format is 2022-01-01 01:00:00
func TSToTime(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

func SubDays(before, last string) int {
	var day int
	t1, _ := time.Parse(TIME_LAYOUT_SECOND, before)
	t2, _ := time.Parse(TIME_LAYOUT_SECOND, last)
	swap := false
	if t1.Unix() > t2.Unix() {
		t1, t2 = t2, t1
		swap = true
	}

	t1_ := t1.Add(time.Duration(t2.Sub(t1).Milliseconds()%86400000) * time.Millisecond)
	day = int(t2.Sub(t1).Hours() / 24)
	// 计算在t1+两个时间的余数之后天数是否有变化
	if t1_.Day() != t1.Day() {
		day += 1
	}

	if swap {
		day = -day
	}

	return day
}

func SubYearSets(before, last time.Time, useFirst bool, useLast bool) []string {
	d := []string{}
	if useFirst {
		d = append(d, before.Format(TIME_LAYOUT_YEAR))
	}

	cursor := before
	for cursor.Before(last) {
		cursor = cursor.AddDate(1, 0, 0)

		if cursor.Before(last) {
			d = append(d, cursor.Format(TIME_LAYOUT_YEAR))
		}
	}

	if !useLast {
		if len(d) > 0 {
			d = d[0 : len(d)-1]
		}
	}

	return d
}

func SubMonSets(before, last time.Time, useFirst bool, useLast bool, seps ...string) []string {
	format := TIME_LAYOUT_MONTH
	if len(seps) > 0 {
		format = strings.ReplaceAll(format, "-", seps[0])
	}

	d := []string{}
	if useFirst {
		d = append(d, before.Format(format))
	}

	cursor := before
	for cursor.Before(last) {
		cursor = cursor.AddDate(0, 1, 0)

		if cursor.Before(last) {
			d = append(d, cursor.Format(format))
		}
	}

	if !useLast {
		if len(d) > 0 {
			d = d[0 : len(d)-1]
		}
	}

	return d
}

func SubDaySets(before, last time.Time, useFirst bool, useLast bool, seps ...string) []string {
	format := TIME_LAYOUT_DAY
	if len(seps) > 0 {
		format = strings.ReplaceAll(format, "-", seps[0])
	}

	d := []string{}
	if useFirst {
		d = append(d, before.Format(format))
	}

	cursor := before
	for cursor.Before(last) {
		cursor = cursor.AddDate(0, 0, 1)

		if cursor.Before(last) {
			d = append(d, cursor.Format(format))
		}
	}

	if !useLast {
		if len(d) > 0 {
			d = d[0 : len(d)-1]
		}
	}

	return d
}
