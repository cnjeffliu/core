/*
 * @Author: Jeffrey.Liu
 * @Date: 2021-07-19 11:58:51
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-12-02 23:28:54
 * @Description:
 */
package logx_test

import (
	"testing"

	"github.com/cnjeffliu/gocore/logx"
)

func TestWriteLog(t *testing.T) {
	logx.Init("output.log")

	logx.Info("init")
	logx.Debugf("%v", "debug info")
}
