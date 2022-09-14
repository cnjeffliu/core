/*
 * @Author: Jeffrey.Liu
 * @Date: 2021-12-15 14:18:20
 * @LastEditors: Jeffrey.Liu
 * @LastEditTime: 2021-12-15 15:40:43
 * @Description:
 */
package thread

import (
	"fmt"

	"github.com/cnjeffliu/gocore/logx"
)

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		// os.Stdout.WriteString(fmt.Sprint(p))
		logx.Error(fmt.Sprint(p))
	}
}
