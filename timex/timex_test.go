/*
 * @Author: Jeffrey Liu
 * @Date: 2022-10-21 23:40:50
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-10-24 21:56:30
 * @Description:
 */
package timex_test

import (
	"testing"
	"time"

	"github.com/cnjeffliu/gocore/timex"
	"github.com/stretchr/testify/assert"
)

func TestSubYearSets(t *testing.T) {
	before := time.Date(2020, 1, 5, 1, 5, 10, 0, time.Local)
	after := time.Date(2021, 5, 6, 13, 10, 30, 0, time.Local)

	useFirst := true
	useLast := true
	assert.Equal(t, []string{"2020", "2021"}, timex.SubYearSets(before, after, useFirst, useLast), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = true
	useLast = false
	assert.Equal(t, []string{"2020"}, timex.SubYearSets(before, after, useFirst, useLast), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = true
	assert.Equal(t, []string{"2021"}, timex.SubYearSets(before, after, useFirst, useLast), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = false
	assert.Equal(t, []string{}, timex.SubYearSets(before, after, useFirst, useLast), "useFirst:%v useLast:%v", useFirst, useLast)
}

func TestSubMonSets(t *testing.T) {
	before := time.Date(2020, 11, 5, 1, 5, 10, 0, time.Local)
	after := time.Date(2021, 1, 6, 13, 10, 30, 0, time.Local)

	useFirst := true
	useLast := true
	assert.Equal(t, []string{"202011", "202012", "202101"}, timex.SubMonSets(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = true
	useLast = false
	assert.Equal(t, []string{"202011", "202012"}, timex.SubMonSets(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = true
	assert.Equal(t, []string{"202012", "202101"}, timex.SubMonSets(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = false
	assert.Equal(t, []string{"202012"}, timex.SubMonSets(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)
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
	assert.Equal(t, []string{"20220130", "20220131", "20220201", "20220202"}, timex.SubDaySets(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = true
	useLast = false
	assert.Equal(t, []string{"20220130", "20220131", "20220201"}, timex.SubDaySets(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = true
	assert.Equal(t, []string{"20220131", "20220201", "20220202"}, timex.SubDaySets(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = false
	assert.Equal(t, []string{"20220131", "20220201"}, timex.SubDaySets(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)
}

func BenchmarkSubDaySets(t *testing.B) {
	before := time.Date(2022, 1, 30, 1, 5, 10, 0, time.Local)
	after := time.Date(2022, 2, 2, 13, 10, 30, 0, time.Local)

	useFirst := true
	useLast := true
	for i := 0; i < t.N; i++ {
		assert.Equal(t, []string{"20220130", "20220131", "20220201", "20220202"}, timex.SubDaySets(before, after, useFirst, useLast, ""), "useFirst:%v useLast:%v", useFirst, useLast)
	}
}
