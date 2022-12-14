/*
 * @Author: Jeffrey Liu
 * @Date: 2022-12-13 12:02:38
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-12-14 16:51:58
 * @Description:
 */
package mathx

import (
	"math/rand"
	"time"
)

// GenerateRandomStr return a random string, size specify the string length.
func GenerateRandomStr(size int) []byte {
	const s = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ"
	var out []byte
	var length = len(s)

	now := time.Now().UnixNano() / 1e6

	basefmt := "2006-01-02 15:04:05"
	t, _ := time.Parse(basefmt, "2020-01-01 01:01:01")
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
