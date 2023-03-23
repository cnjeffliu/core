/*
 * @Author: cnzf1
 * @Date: 2021-07-19 11:58:51
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-17 18:37:14
 * @Description:
 */
package logx_test

import (
	"testing"

	"github.com/cnzf1/gocore/logx"
)

func TestWriteLog(t *testing.T) {
	logx.Init(logx.WithPath("output.log"))

	logx.Info("init")
	logx.Debugf("%v", "debug info")
}
