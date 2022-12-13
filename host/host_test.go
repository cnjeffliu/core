/*
 * @Author: Jeffrey Liu
 * @Date: 2022-09-13 20:37:28
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-12-13 14:11:47
 * @Description:
 */
package host_test

import (
	"fmt"
	"testing"

	"github.com/cnzf1/gocore/host"
	"github.com/cnzf1/gocore/timex"
)

func TestBTime(t *testing.T) {
	tt := host.GetBtime()
	fmt.Println(tt)

	fmt.Println(timex.TSToStr(tt, 0))
}
