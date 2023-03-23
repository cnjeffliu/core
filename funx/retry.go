package funx

import (
	"time"

	"github.com/cnzf1/gocore/errorx"
)

const defaultRetryTimes = 3
const defaultPeriod = 0

type (
	// RetryOption defines the method to customize DoWithRetry.
	RetryOption func(*retryOptions)

	retryOptions struct {
		times  int
		period time.Duration
	}
)

// DoWithRetry runs fn, and retries if failed. Default to retry 3 times.
func DoWithRetry(fn func() error, opts ...RetryOption) error {
	options := newRetryOptions()
	for _, opt := range opts {
		opt(options)
	}

	var berr errorx.BatchError
	for i := 0; i < options.times; i++ {
		if err := fn(); err != nil {
			time.Sleep(options.period)
			berr.Add(err)
		} else {
			return nil
		}
	}

	return berr.Err()
}

// WithRetry customize a DoWithRetry call with given retry times.
func WithRetry(times int) RetryOption {
	return func(options *retryOptions) {
		options.times = times
	}
}

// WithPeriod customize a DoWithRetry call with given period(ms).
func WithPeriod(period time.Duration) RetryOption {
	return func(options *retryOptions) {
		options.period = period
	}
}

func newRetryOptions() *retryOptions {
	return &retryOptions{
		times:  defaultRetryTimes,
		period: defaultPeriod,
	}
}
