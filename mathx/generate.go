/*
 * @Author: cnzf1
 * @Date: 2022-12-13 12:02:38
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-29 11:23:09
 * @Description:
 */
package mathx

import (
	"math/rand"
	"time"

	"github.com/cnzf1/gocore/timex"
)

// RandStr return a random string, size specify the string length.
func RandStr(size int) []byte {
	const s = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ"
	var out []byte
	var length = len(s)

	now := timex.NowMs()

	t, _ := time.Parse(timex.TIME_LAYOUT_SECOND, "2020-01-01 01:01:01")
	base := t.UnixNano() / 1e6

	period := now - base
	for period > 0 {
		rest := period % int64(length)
		period = (period - rest) / int64(length)
		out = append(out, s[rest])
	}

	if len(out) >= size {
		return out[:size]
	}

	for len(out) < size {
		out = append(out, s[rand.Intn(length)])
	}
	return out
}
