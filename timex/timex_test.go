/*
 * @Author: Jeffrey Liu
 * @Date: 2022-10-21 23:40:50
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-10-22 00:31:02
 * @Description:
 */
package timex_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cnjeffliu/gocore/timex"
)

func TestSubYearSets(t *testing.T) {
	before := time.Date(2020, 1, 5, 1, 5, 10, 0, time.Local)
	after := time.Date(2021, 5, 6, 13, 10, 30, 0, time.Local)
	fmt.Println("before:", before)
	fmt.Println("after:", after)

	useFirst := true
	useLast := true
	fmt.Println("useFirst:", useFirst, " useLast:", useLast, " ", timex.SubYearSets(before, after, useFirst, useLast).List())

	useFirst = false
	useLast = false
	fmt.Println("useFirst:", useFirst, " useLast:", useLast, " ", timex.SubYearSets(before, after, useFirst, useLast).List())

	useFirst = true
	useLast = false
	fmt.Println("useFirst:", useFirst, " useLast:", useLast, " ", timex.SubYearSets(before, after, useFirst, useLast).List())

	useFirst = false
	useLast = true
	fmt.Println("useFirst:", useFirst, " useLast:", useLast, " ", timex.SubYearSets(before, after, useFirst, useLast).List())

}

func TestSubMonSets(t *testing.T) {
	before := time.Date(2020, 11, 5, 1, 5, 10, 0, time.Local)
	after := time.Date(2021, 1, 6, 13, 10, 30, 0, time.Local)
	fmt.Println("before:", before)
	fmt.Println("after:", after)

	useFirst := true
	useLast := true
	fmt.Println("useFirst:", useFirst, " useLast:", useLast, " ", timex.SubMonSets(before, after, useFirst, useLast, "-").List())

	useFirst = false
	useLast = false
	fmt.Println("useFirst:", useFirst, " useLast:", useLast, " ", timex.SubMonSets(before, after, useFirst, useLast, "-").List())

	useFirst = true
	useLast = false
	fmt.Println("useFirst:", useFirst, " useLast:", useLast, " ", timex.SubMonSets(before, after, useFirst, useLast, "-").List())

	useFirst = false
	useLast = true
	fmt.Println("useFirst:", useFirst, " useLast:", useLast, " ", timex.SubMonSets(before, after, useFirst, useLast, "-").List())
}

func TestSubDaySets(t *testing.T) {
	before := time.Date(2022, 1, 30, 1, 5, 10, 0, time.Local)
	after := time.Date(2022, 2, 2, 13, 10, 30, 0, time.Local)
	fmt.Println("before:", before)
	fmt.Println("after:", after)

	useFirst := true
	useLast := true
	fmt.Println("useFirst:", useFirst, " useLast:", useLast, " ", timex.SubDaySets(before, after, useFirst, useLast, "-").List())

	useFirst = false
	useLast = false
	fmt.Println("useFirst:", useFirst, " useLast:", useLast, " ", timex.SubDaySets(before, after, useFirst, useLast, "-").List())

	useFirst = true
	useLast = false
	fmt.Println("useFirst:", useFirst, " useLast:", useLast, " ", timex.SubDaySets(before, after, useFirst, useLast, "-").List())

	useFirst = false
	useLast = true
	fmt.Println("useFirst:", useFirst, " useLast:", useLast, " ", timex.SubDaySets(before, after, useFirst, useLast, "-").List())
}
