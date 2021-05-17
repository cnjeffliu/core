package logx

import "testing"

func TestWriteLog(t *testing.T) {
	InitLog("output.log")

	Info("init")
	Debug("test")
}
