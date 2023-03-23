/*
 * @Author: cnzf1
 * @Date: 2023-03-28 14:36:32
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-29 13:47:13
 * @Description:
 */
package mapx

import (
	"sync"
	"time"

	"github.com/cnzf1/gocore/lang"
	"github.com/cnzf1/gocore/thread"
	"go.uber.org/atomic"
)

type ExpiredMap struct {
	m        sync.Map
	tick     time.Duration
	capacity int64
	size     atomic.Int64
	stop     chan bool
	delFn    DelCallBack
}

type mapItem struct {
	key    string
	value  lang.AnyType
	expire atomic.Time
}

type EMConfig struct {
	tick  time.Duration
	delFn DelCallBack
}

type EMOption func(*EMConfig)

func WithTick(tick time.Duration) EMOption {
	return func(e *EMConfig) {
		e.tick = tick
	}
}

type DelCallBack func(key string, val lang.AnyType)

func WithDelCallback(fn DelCallBack) EMOption {
	return func(e *EMConfig) {
		e.delFn = fn
	}
}
func NewExpiredMap(opts ...EMOption) *ExpiredMap {
	cfg := &EMConfig{
		tick:  time.Second,
		delFn: nil,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	c := &ExpiredMap{
		m:     sync.Map{},
		tick:  cfg.tick,
		stop:  make(chan bool),
		delFn: cfg.delFn,
	}

	c.check()
	return c
}

func (c *ExpiredMap) check() {
	thread.GoSafe(func() {
		ticker := time.NewTicker(c.tick)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.m.Range(func(key, value any) bool {
					v, ok := value.(*mapItem)
					if ok && v != nil && time.Now().After(v.expire.Load()) {
						tkey := key.(string)
						c.Delete(tkey)
					}
					return true
				})
			case <-c.stop:
				return
			}
		}
	})
}

func (c *ExpiredMap) Set(key string, value lang.AnyType, ttl time.Duration) bool {
	v := &mapItem{
		key:   key,
		value: value,
	}

	v.expire.Store(time.Now().Add(ttl))
	c.m.Store(key, v)
	c.size.Inc()
	return true
}

func (c *ExpiredMap) Get(key string) (value lang.AnyType, ok bool) {
	var v any
	v, ok = c.m.Load(key)
	if !ok {
		return
	}

	v2, ok := v.(*mapItem)
	if !ok || v2 == nil {
		return
	}

	if time.Now().After(v2.expire.Load()) {
		ok = false
		c.Delete(key)
		return
	}

	value = v2.value
	return
}

func (c *ExpiredMap) Delete(key string) {
	val, _ := c.Get(key)
	c.m.Delete(key)
	c.size.Dec()
	if c.delFn != nil {
		c.delFn(key, val)
	}
}

func (c *ExpiredMap) Size() int64 {
	return c.size.Load()
}

func (c *ExpiredMap) TTL(key string) time.Duration {
	v, ok := c.m.Load(key)
	if !ok {
		return -2
	}

	v2, ok := v.(*mapItem)
	if !ok || v2 == nil {
		return -2
	}

	now := time.Now()
	if now.After(v2.expire.Load()) {
		c.Delete(key)
		return -1
	}

	return v2.expire.Load().Sub(now)
}

func (c *ExpiredMap) Foreach(fn func(key string, value lang.AnyType)) {
	c.m.Range(func(key, value any) bool {
		v, ok := value.(*mapItem)
		if !ok || v == nil {
			return false
		}

		if time.Now().After(v.expire.Load()) {
			tkey := key.(string)
			c.Delete(tkey)
			return true
		}
		tkey := key.(string)
		fn(tkey, v.value)
		return true
	})
}

func (c *ExpiredMap) Clear() {
	c.m.Range(func(key, value any) bool {
		tkey := key.(string)
		c.Delete(tkey)
		return true
	})
}

func (c *ExpiredMap) Close() {
	c.stop <- true
	c.Clear()
}
