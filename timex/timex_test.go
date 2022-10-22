/*
 * @Author: Jeffrey Liu
 * @Date: 2022-10-21 23:40:50
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-10-22 18:07:06
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
	assert.Equal(t, []string{"2020", "2021"}, timex.SubYearSets(before, after, useFirst, useLast).List(), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = true
	useLast = false
	assert.Equal(t, []string{"2020"}, timex.SubYearSets(before, after, useFirst, useLast).List(), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = true
	assert.Equal(t, []string{"2021"}, timex.SubYearSets(before, after, useFirst, useLast).List(), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = false
	assert.Equal(t, []string{}, timex.SubYearSets(before, after, useFirst, useLast).List(), "useFirst:%v useLast:%v", useFirst, useLast)
}

func TestSubMonSets(t *testing.T) {
	before := time.Date(2020, 11, 5, 1, 5, 10, 0, time.Local)
	after := time.Date(2021, 1, 6, 13, 10, 30, 0, time.Local)

	useFirst := true
	useLast := true
	assert.Equal(t, []string{"202011", "202012", "202101"}, timex.SubMonSets(before, after, useFirst, useLast).List(), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = true
	useLast = false
	assert.Equal(t, []string{"202011", "202012"}, timex.SubMonSets(before, after, useFirst, useLast).List(), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = true
	assert.Equal(t, []string{"202012", "202101"}, timex.SubMonSets(before, after, useFirst, useLast).List(), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = false
	assert.Equal(t, []string{"202012"}, timex.SubMonSets(before, after, useFirst, useLast).List(), "useFirst:%v useLast:%v", useFirst, useLast)
}

func TestSubDaySets(t *testing.T) {
	before := time.Date(2022, 1, 30, 1, 5, 10, 0, time.Local)
	after := time.Date(2022, 2, 2, 13, 10, 30, 0, time.Local)

	useFirst := true
	useLast := true
	assert.Equal(t, []string{"20220130", "20220131", "20220201", "20220202"}, timex.SubDaySets(before, after, useFirst, useLast).List(), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = true
	useLast = false
	assert.Equal(t, []string{"20220130", "20220131", "20220201"}, timex.SubDaySets(before, after, useFirst, useLast).List(), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = true
	assert.Equal(t, []string{"20220131", "20220201", "20220202"}, timex.SubDaySets(before, after, useFirst, useLast).List(), "useFirst:%v useLast:%v", useFirst, useLast)

	useFirst = false
	useLast = false
	assert.Equal(t, []string{"20220131", "20220201"}, timex.SubDaySets(before, after, useFirst, useLast).List(), "useFirst:%v useLast:%v", useFirst, useLast)
}
