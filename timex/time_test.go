/*
 * @Author: cnzf1
 * @Date: 2022-10-21 23:40:50
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-25 23:39:20
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

func TestSubYearSets(t *testing.T) {
	before := time.Date(2020, 1, 5, 1, 5, 10, 0, time.Local)
	after := time.Date(2021, 5, 6, 13, 10, 30, 0, time.Local)

	useFirst := true
	useLast := true
	assert.Equal(t, []string{"2020", "2021"}, timex.SubYearSetsEx(before, after, useFirst, useLast), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = true
	useLast = false
	assert.Equal(t, []string{"2020"}, timex.SubYearSetsEx(before, after, useFirst, useLast), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = true
	assert.Equal(t, []string{"2021"}, timex.SubYearSetsEx(before, after, useFirst, useLast), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = false
	assert.Equal(t, []string{}, timex.SubYearSetsEx(before, after, useFirst, useLast), "useFirst:%v useLast:%v", useFirst, useLast)
}

func TestSubMonSets(t *testing.T) {
	before := time.Date(2020, 11, 5, 1, 5, 10, 0, time.Local)
	after := time.Date(2021, 1, 6, 13, 10, 30, 0, time.Local)

	useFirst := true
	useLast := true
	assert.Equal(t, []string{"202011", "202012", "202101"}, timex.SubMonSetsEx(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = true
	useLast = false
	assert.Equal(t, []string{"202011", "202012"}, timex.SubMonSetsEx(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = true
	assert.Equal(t, []string{"202012", "202101"}, timex.SubMonSetsEx(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = false
	assert.Equal(t, []string{"202012"}, timex.SubMonSetsEx(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)
}

func TestSubDays(t *testing.T) {
	before := time.Date(2020, 12, 28, 1, 5, 10, 0, time.Local)
	after := time.Date(2021, 1, 2, 13, 10, 30, 0, time.Local)

	assert.Equal(t, 5, timex.SubDays(before, after))

	before = time.Date(2021, 1, 1, 23, 59, 59, 0, time.Local)
	after = time.Date(2021, 1, 2, 0, 0, 1, 0, time.Local)

	assert.Equal(t, 1, timex.SubDays(before, after))
}

func TestSubDaySets(t *testing.T) {
	before := time.Date(2022, 1, 30, 1, 5, 10, 0, time.Local)
	after := time.Date(2022, 2, 2, 13, 10, 30, 0, time.Local)

	useFirst := true
	useLast := true
	assert.Equal(t, []string{"20220130", "20220131", "20220201", "20220202"}, timex.SubDaySetsEx(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = true
	useLast = false
	assert.Equal(t, []string{"20220130", "20220131", "20220201"}, timex.SubDaySetsEx(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = true
	assert.Equal(t, []string{"20220131", "20220201", "20220202"}, timex.SubDaySetsEx(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = false
	assert.Equal(t, []string{"20220131", "20220201"}, timex.SubDaySetsEx(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)
}

func BenchmarkSubDaySets(t *testing.B) {
	before := time.Date(2022, 1, 30, 1, 5, 10, 0, time.Local)
	after := time.Date(2022, 2, 2, 13, 10, 30, 0, time.Local)

	useFirst := true
	useLast := true
	for i := 0; i < t.N; i++ {
		assert.Equal(t, []string{"20220130", "20220131", "20220201", "20220202"}, timex.SubDaySetsEx(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)
	}
}

func TestNowStr(t *testing.T) {
	type args struct {
		layout string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: timex.TIME_LAYOUT_SECOND,
			args: args{
				layout: timex.TIME_LAYOUT_SECOND,
			},
			want: time.Now().Format(timex.TIME_LAYOUT_SECOND),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timex.NowStr(tt.args.layout); got != tt.want {
				t.Errorf("NowStr() = %v, want %v", got, tt.want)
			}
		})
	}

	fmt.Println(timex.NowStr(timex.TIME_LAYOUT_COMPACT_MILLSECOND))
}
