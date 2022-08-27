/*
 * @Author: Jeffrey Liu <zhifeng172@163.com>
 * @Date: 2022-07-20 13:56:45
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-08-26 11:43:42
 * @Description:
 */
package timex

import "time"

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
