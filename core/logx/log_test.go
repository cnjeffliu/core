/*
 * @Author: Jeffrey.Liu <zhifeng172@163.com>
 * @Date: 2021-07-19 11:58:51
 * @LastEditors: Jeffrey.Liu
 * @LastEditTime: 2021-12-15 16:25:35
 * @Description:
 */
package logx

import "testing"

func TestWriteLog(t *testing.T) {
	Init(WithFile("output.log"))

	Info("init")
	Debug("test")
}
