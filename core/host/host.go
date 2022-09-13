/*
 * @Author: Jeffrey Liu <zhifeng172@163.com>
 * @Date: 2022-09-13 20:32:00
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-09-13 20:42:45
 * @Description:
 */
package host

import (
	"strconv"
	"strings"

	"gitee.com/cnjeffliu/core/filex"
)

func GetBtime() int64 {
	file := "/proc/stat"
	lines, err := filex.ReadLinesOffsetN(file, 0, -1)
	if err != nil {
		return 0
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "btime") {
			f := strings.Fields(line)
			if len(f) != 2 {
				return 0
			}

			b, err := strconv.ParseInt(f[1], 10, 64)
			if err != nil {
				return 0
			}

			t := int64(b)
			return t
		}
	}

	return 0
}
