package limit_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/cnzf1/gocore/collection/limit"
)

func TestWindowLimit(t *testing.T) {
	wl := limit.NewWindowLimit(limit.WithLimit(500), limit.WithPeriod(1))
	total := 0
	succ := 0

	for i := 0; i < 10000; i++ {
		ok := wl.Access()
		total += 1
		if ok {
			succ += 1
		}

		// time.Sleep(5 * time.Millisecond)
	}

	// for i := 0; i < 10; i++ {
	// 	fmt.Println("cur items:", wl.Count())
	// 	time.Sleep(100 * time.Microsecond)
	// }

	fmt.Println("total:", total, " succ:", succ)
}

func TestMultWindowLimit(t *testing.T) {
	keys := [...]string{"10.0.2.93", "10.0.2.113"}
	var sm sync.Map
	var wg sync.WaitGroup

	for _, val := range keys {
		key := val
		wg.Add(1)

		wlp, ok := sm.Load(key)
		if !ok {
			sm.Store(key, limit.NewWindowLimit(limit.WithLimit(500), limit.WithPeriod(10)))
			wlp, _ = sm.Load(key)
		}
		wl := wlp.(*limit.WindowLimit)

		go func() {
			total := 0
			succ := 0

			for i := 0; i < 10000; i++ {
				ok := wl.Access()
				total += 1
				if ok {
					succ += 1
				}

				time.Sleep((time.Duration(1 + rand.Intn(10))) * time.Millisecond)
			}

			t.Log("key:", key, "total:", total, " succ:", succ)
			wg.Done()
		}()
	}

	wg.Wait()
}
