/**
 * @Author: Jeffrey.Liu
 * @Date: 2021-11-18 18:46:30
 * @LastEditors: Jeffrey.Liu
 * @LastEditTime: 2021-11-18 18:46:30
 * @Description:
 */

package demo

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := StatusOK.WithErr(errors.New("with err msg"))
	fmt.Println(String(err, ErrLangEN))
}
