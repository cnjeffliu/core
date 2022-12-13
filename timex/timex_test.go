/*
 * @Author: Jeffrey Liu
 * @Date: 2022-10-21 23:40:50
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-12-01 18:26:59
 * @Description:
 */
package timex_test

import (
	"testing"
	"time"

	"github.com/cnzf1/gocore/timex"
	"github.com/stretchr/testify/assert"
)

func TestSubYearSets(t *testing.T) {
	before := time.Date(2020, 1, 5, 1, 5, 10, 0, time.Local)
	after := time.Date(2021, 5, 6, 13, 10, 30, 0, time.Local)

	list := timex.SubYearSets(before, after, true, true)
	assert.Equal(t, []string{"2020", "2021"}, list)

	list = timex.SubYearSets(before, after, true, false)
	assert.Equal(t, []string{"2020"}, list)

	list = timex.SubYearSets(before, after, false, true)
	assert.Equal(t, []string{"2021"}, list)

	list = timex.SubYearSets(before, after, false, false)
	assert.Equal(t, []string{}, list)
}

func TestSubMonSets(t *testing.T) {
	before := time.Date(2020, 11, 5, 1, 5, 10, 0, time.Local)
	after := time.Date(2021, 1, 6, 13, 10, 30, 0, time.Local)

	list := timex.SubMonSets(before, after, true, true, "")
	assert.Equal(t, []string{"202011", "202012", "202101"}, list)

	list = timex.SubMonSets(before, after, true, false, "")
	assert.Equal(t, []string{"202011", "202012"}, list)

	list = timex.SubMonSets(before, after, false, true, "")
	assert.Equal(t, []string{"202012", "202101"}, list)

	list = timex.SubMonSets(before, after, false, false, "")
	assert.Equal(t, []string{"202012"}, list)
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

	list := timex.SubDaySets(before, after, true, true, "")
	assert.Equal(t, []string{"20220130", "20220131", "20220201", "20220202"}, list)

	list = timex.SubDaySets(before, after, true, false, "")
	assert.Equal(t, []string{"20220130", "20220131", "20220201"}, list)

	list = timex.SubDaySets(before, after, false, true, "")
	assert.Equal(t, []string{"20220131", "20220201", "20220202"}, list)

	list = timex.SubDaySets(before, after, false, false, "")
	assert.Equal(t, []string{"20220131", "20220201"}, list)
}

func TestSubDaySets2(t *testing.T) {
	start := "2022-11-29"
	end := "2022-12-01"
	before := timex.DateStrToTime(start)
	after := timex.DateStrToTime(end)

	list := timex.SubDaySets(before, after, true, true, "")
	assert.Equal(t, []string{"20221129", "20221130"}, list)
}

func BenchmarkSubDaySets(t *testing.B) {
	before := time.Date(2022, 1, 30, 1, 5, 10, 0, time.Local)
	after := time.Date(2022, 2, 2, 13, 10, 30, 0, time.Local)

	for i := 0; i < t.N; i++ {
		assert.Equal(t, []string{"20220130", "20220131", "20220201", "20220202"}, timex.SubDaySets(before, after, true, true, ""))
	}
}
