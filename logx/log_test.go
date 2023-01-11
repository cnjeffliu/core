/*
 * @Author: cnzf1
 * @Date: 2021-07-19 11:58:51
 * @LastEditors: cnzf1
 * @LastEditTime: 2022-12-02 23:28:54
 * @Description:
 */
package logx_test

import (
	"testing"

	"github.com/cnzf1/gocore/logx"
)

func TestWriteLog(t *testing.T) {
	logx.Init("output.log")

	logx.Info("init")
	logx.Debugf("%v", "debug info")
}
