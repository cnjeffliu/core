/*
 * @Author: cnzf1
 * @Date: 2021-08-18 18:07:12
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-27 13:43:26
 * @Description: 控制超时时间运行函数
 */
package funx

import (
	"context"
	"time"

	"github.com/cnzf1/gocore/lang"
	"github.com/cnzf1/gocore/thread"
)

var (
	// ErrCanceled is the error returned when the context is canceled.
	ErrCanceled = context.Canceled
	// ErrTimeout is the error returned when the context's deadline passes.
	ErrTimeout = context.DeadlineExceeded
)

type DoConfig struct {
	ctx     context.Context
	timeout time.Duration
	exit    <-chan lang.PlaceholderType
}

// DoOption defines the method to customize a DoWithTimeout call.
type DoOption func(c *DoConfig)

// WithContext customizes a DoWithTimeout call with given ctx.
func WithContext(ctx context.Context) DoOption {
	return func(c *DoConfig) {
		c.ctx = ctx
	}
}

// WithStopChan
func WithStopChan(ch <-chan lang.PlaceholderType) DoOption {
	return func(c *DoConfig) {
		c.exit = ch
	}
}

func WithTimeout(tm time.Duration) DoOption {
	return func(c *DoConfig) {
		c.timeout = tm
	}
}

// DoWithTimeout runs fn with timeout control.
func DoWithTimeout(fn func() error, opts ...DoOption) error {
	cfg := &DoConfig{
		ctx:     context.Background(),
		timeout: 5 * time.Second,
		exit:    make(<-chan lang.PlaceholderType),
	}
	for _, opt := range opts {
		opt(cfg)
	}
	done := make(chan error, 1)

	ctx, cancel := context.WithTimeout(cfg.ctx, cfg.timeout)
	defer cancel()

	thread.GoSafe(func() {
		done <- fn()
	})

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-cfg.exit:
		return nil
	}
}
