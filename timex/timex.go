/*
 * @Author: Jeffrey Liu
 * @Date: 2022-07-20 13:56:45
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-12-15 14:55:06
 * @Description:
 */
package timex

import (
	"strings"
	"time"
)

const (
	TIME_LAYOUT_YEAR                  = "2006"                          // YYYY
	TIME_LAYOUT_MONTH                 = "2006-01"                       // YYYY-MM
	TIME_LAYOUT_DAY                   = "2006-01-02"                    // YYYY-MM-DD
	TIME_LAYOUT_MINUTE                = "2006-01-02 15:04"              // YYYY-MM-DD HH:mm
	TIME_LAYOUT_SECOND                = "2006-01-02 15:04:05"           // YYYY-MM-DD HH-mm:SS
	TIME_LAYOUT_MILLSECOND            = "2006-01-02 15:04:05.000"       // YYYY-MM-DD HH-mm:SS.NNN
	TIME_LAYOUT_MICROSECOND           = "2006-01-02 15:04:05.000000"    // YYYY-MM-DD HH-mm:SS.NNNNNN
	TIME_LAYOUT_NANOSECOND            = "2006-01-02 15:04:05.000000000" // YYYY-MM-DD HH-mm:SS.NNNNNNNNN
	TIME_LAYOUT_COMPACT_DAY           = "20060102"                      // YYYYMMDD
	TIME_LAYOUT_COMPACT_SIMPLE_SECOND = "150405"                        // hhmmss
	TIME_LAYOUT_COMPACT_EXT_SECOND    = "0102150405"                    // MMDDhhmmss
	TIME_LAYOUT_COMPACT_SECOND        = "20060102150405"                // YYYYMMDDhhmmss
)

// NowSecond returns the current seconds.
func NowSecond() int64 {
	return time.Now().Unix()
}

// NowMillis returns the current milliseconds.
func NowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// NowMicros returns the current microseconds.
func NowMicros() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

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

// SubDays return the days between before and last
func SubDays(before, last time.Time) int {
	var day int
	swap := false
	if before.Unix() > last.Unix() {
		before, last = last, before
		swap = true
	}

	t1_ := before.Add(time.Duration(last.Sub(before).Milliseconds()%86400000) * time.Millisecond)
	day = int(last.Sub(before).Hours() / 24)
	// 计算在t1+两个时间的余数之后天数是否有变化
	if t1_.Day() != before.Day() {
		day += 1
	}

	if swap {
		day = -day
	}

	return day
}

// SubYearSets return the year set during before and last.
// useFirst represent include the year of before.
// useLast represent include the year of last.
//
// For example:
//
//	before := time.Date(2020, 12, 28, 1, 5, 10, 0, time.Local)
//	after := time.Date(2021, 1, 2, 13, 10, 30, 0, time.Local)
//	Call SubYearSets(before, after, true, true)
//		-> []string{"2020", "2021"}
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

// SubMonSets return the month set during before and last.
// useFirst represent include the month of before.
// useLast represent include the month of last.
//
// For example:
//
//	before := time.Date(2020, 12, 28, 1, 5, 10, 0, time.Local)
//	after := time.Date(2021, 1, 2, 13, 10, 30, 0, time.Local)
//	Call SubMonSets(before, after, true, true)
//		-> []string{"2020-12", "2021-01"}
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

// SubDaySets return the date set during before and last.
// useFirst represent include the date of before.
// useLast represent include the date of last.
//
// For example:
//
//	before := time.Date(2020, 12, 28, 1, 5, 10, 0, time.Local)
//	after := time.Date(2021, 1, 2, 13, 10, 30, 0, time.Local)
//	Call SubDaySets(before, after, true, true)
//		-> []string{"2020-12-28","2020-12-29","2020-12-30", "2020-12-31","2021-01-01","2021-01-02"}
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
