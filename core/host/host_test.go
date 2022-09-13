/*
 * @Author: Jeffrey Liu <zhifeng172@163.com>
 * @Date: 2022-09-13 20:37:28
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-09-13 20:42:27
 * @Description:
 */
package host_test

import (
	"fmt"
	"testing"

	"gitee.com/cnjeffliu/core/host"
	"gitee.com/cnjeffliu/core/timex"
)

func TestBTime(t *testing.T) {
	tt := host.GetBtime()
	fmt.Println(tt)

	fmt.Println(timex.TSToStr(tt, 0))
}
