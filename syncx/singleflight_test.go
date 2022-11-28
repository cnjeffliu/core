/*
 * @Author: Jeffrey Liu
 * @Date: 2022-11-25 16:16:59
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-11-28 14:33:09
 * @Description:
 */
package syncx

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestExclusiveCallDoDiffDupSuppress(t *testing.T) {
	g := NewSingleFlight()
	broadcast := make(chan struct{})
	var calls int32
	tests := []string{"e", "a", "e", "a", "b", "c", "b", "a", "c", "d", "b", "c", "d"}

	var wg sync.WaitGroup
	for _, key := range tests {
		wg.Add(1)
		go func(k string) {
			<-broadcast // get all goroutines ready
			_, err := g.Do(k, func() (interface{}, error) {
				atomic.AddInt32(&calls, 1)
				time.Sleep(10 * time.Millisecond)
				return nil, nil
			})
			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			wg.Done()
		}(key)
	}

	time.Sleep(100 * time.Millisecond)
	close(broadcast)
	wg.Wait()

	if got := atomic.LoadInt32(&calls); got != 5 {
		// five letters
		t.Errorf("number of calls = %d; want 5", got)
	}
}

func TestExclusiveCallDoExDupSuppress(t *testing.T) {
	g := NewSingleFlight()
	c := make(chan string)
	var calls int32
	fn := func() (interface{}, error) {
		atomic.AddInt32(&calls, 1)
		return <-c, nil
	}

	const n = 10
	var wg sync.WaitGroup
	var freshes int32
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			v, fresh, err := g.DoEx("key", fn)
			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			if fresh {
				atomic.AddInt32(&freshes, 1)
			}
			if v.(string) != "bar" {
				t.Errorf("got %q; want %q", v, "bar")
			}
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond) // let goroutines above block
	c <- "bar"
	wg.Wait()
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("number of calls = %d; want 1", got)
	}
	if got := atomic.LoadInt32(&freshes); got != 1 {
		t.Errorf("freshes = %d; want 1", got)
	}
}
