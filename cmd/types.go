/*
 * @Author: Jeffrey <zhifeng172@163.com>
 * @Date: 2021-07-29 20:11:51
 * @Descripttion:
 */
package cmd

import (
	"serv/core/errorx"
)

const (
	// BadRequest error
	BadRequest errorx.ErrorType = iota + 1
	// NotFound error
	NotFound
)
