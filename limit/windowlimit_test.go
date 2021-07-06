package limit

import (
	"fmt"
	"testing"
	"time"
)

func TestWindowLimit(t *testing.T) {
	wl := NewWindowLimit(500, 1)
	total := 0
	succ := 0

	for i := 0; i < 10000; i++ {
		ok := wl.Access()
		total += 1
		if ok {
			succ += 1
		}

		time.Sleep(5 * time.Millisecond)
	}

	// for i := 0; i < 10; i++ {
	// 	fmt.Println("cur items:", wl.Count())
	// 	time.Sleep(100 * time.Microsecond)
	// }

	fmt.Println("total:", total, " succ:", succ)
}
