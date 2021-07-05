package limit

import (
	"fmt"
	"testing"
)

func TestSlidingWinSingle(t *testing.T) {
	var total uint64
	var succ uint64
	for i := 0; i < 10000; i++ {
		ok := SlidingWinSingle("key", 500, 1)
		// fmt.Printf("%v idx:%v ok? %v \n", time.Now().Unix(), i, ok)
		// time.Sleep(time.Millisecond * 100)
		total += 1
		if ok {
			succ += 1
		}
	}

	fmt.Println("total:", total, " succ:", succ)
}
