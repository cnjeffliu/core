package thread

import (
	"fmt"
	"serv/core/logx"
)

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		logx.Fatal(fmt.Sprint(p))
	}
}
