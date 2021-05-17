package util

import "testing"

func TestWriteLog(t *testing.T) {
	InitLog("output.log")

	Debug("test")
}
