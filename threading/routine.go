package threading

import (
	"bytes"
	"runtime"
	"strconv"

	"serv/rescue"
)

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RoutineId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]

	n, _ := strconv.ParseUint(string(b), 10, 64)

	return n
}

func RunSafe(fn func()) {
	defer rescue.Recover()

	fn()
}
