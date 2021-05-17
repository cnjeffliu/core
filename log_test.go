package util

import "testing"

func TestWriteLog(t *testing.T) {
	InitLog("D:\\output.log")

	Debug("test")
}
