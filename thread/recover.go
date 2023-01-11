/*
 * @Author: cnzf1
 * @Date: 2021-12-15 14:18:20
 * @LastEditors: cnzf1
 * @LastEditTime: 2022-12-02 22:07:45
 * @Description:
 */
package thread

import (
	"errors"
	"fmt"
	"os"
)

func Recover(cleanups ...func()) (err error) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		s := fmt.Sprint(p)
		os.Stdout.WriteString(s)
		return errors.New(s)
	}
	return nil
}
