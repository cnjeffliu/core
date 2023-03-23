/*
 * @Author: cnzf1
 * @Date: 2023-03-28 15:19:42
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-29 11:34:37
 * @Description:
 */
package mapx_test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/cnzf1/gocore/collection/mapx"
)

type cacheItem struct {
	itemName  string
	itemValue int
}

func TestNewExpiredMap(t *testing.T) {
	cnt := 10000 * 100
	em := mapx.NewExpiredMap()
	wg := &sync.WaitGroup{}

	for i := 0; i < cnt; i++ {
		oneI := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			setRes := em.Set(strconv.Itoa(oneI), &cacheItem{
				itemName:  fmt.Sprintf("new item:%d", oneI),
				itemValue: oneI,
			}, time.Second*60)
			if !setRes {
				t.Errorf("Set error")
				return
			}
		}()
	}

	for i := 0; i < cnt; i++ {
		oneI := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = em.Get(strconv.Itoa(oneI))
		}()
	}

	wg.Wait()
}
