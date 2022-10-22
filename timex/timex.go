/*
 * @Author: Jeffrey Liu
 * @Date: 2022-07-20 13:56:45
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-10-22 18:03:00
 * @Description:
 */
package timex

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cnjeffliu/gocore/setx"
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

func SubYearSets(before, last time.Time, useFirst bool, useLast bool) setx.String {
	tmp := ""
	set := setx.NewString()
	if useFirst {
		tmp = strconv.Itoa(before.Year())
		set.Insert(tmp)
	}

	cursor := before
	lastone := ""
	for cursor.Before(last) {
		cursor = cursor.AddDate(1, 0, 0)
		tmp = strconv.Itoa(cursor.Year())

		if cursor.Before(last) {
			set.Insert(tmp)
			lastone = tmp
		}
	}

	if len(lastone) > 0 && !useLast {
		set.Delete(lastone)
	}

	return set
}

func SubMonSets(before, last time.Time, useFirst bool, useLast bool, seps ...string) setx.String {
	sep := ""
	if len(seps) > 0 {
		sep = seps[0]
	}

	tmp := ""
	set := setx.NewString()
	if useFirst {
		tmp = fmt.Sprintf("%d%s%02d", before.Year(), sep, before.Month())
		set.Insert(tmp)
	}

	cursor := before
	lastone := ""
	for cursor.Before(last) {
		cursor = cursor.AddDate(0, 1, 0)
		tmp = fmt.Sprintf("%d%s%02d", cursor.Year(), sep, cursor.Month())

		if cursor.Before(last) {
			set.Insert(tmp)
			lastone = tmp
		}
	}

	if len(lastone) > 0 && !useLast {
		set.Delete(lastone)
	}

	return set
}

func SubDaySets(before, last time.Time, useFirst bool, useLast bool, seps ...string) setx.String {
	sep := ""
	if len(seps) > 0 {
		sep = seps[0]
	}

	tmp := ""
	set := setx.NewString()
	if useFirst {
		tmp = fmt.Sprintf("%d%s%02d%s%02d", before.Year(), sep, before.Month(), sep, before.Day())
		set.Insert(tmp)
	}

	cursor := before
	lastone := ""
	for cursor.Before(last) {
		cursor = cursor.AddDate(0, 0, 1)
		tmp = fmt.Sprintf("%d%s%02d%s%02d", cursor.Year(), sep, cursor.Month(), sep, cursor.Day())

		if cursor.Before(last) {
			set.Insert(tmp)
			lastone = tmp
		}
	}

	if len(lastone) > 0 && !useLast {
		set.Delete(lastone)
	}

	return set
}
