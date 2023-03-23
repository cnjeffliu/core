package funx_test

import (
	"context"
	"testing"
	"time"

	"github.com/cnzf1/gocore/funx"
	"github.com/stretchr/testify/assert"
)

func TestDoWithTimeout(t *testing.T) {
	err := funx.DoWithTimeout(func() error {
		return nil
	}, funx.WithTimeout(time.Second))
	assert.Nil(t, err)

	err = funx.DoWithTimeout(func() error {
		time.Sleep(time.Millisecond * 100)
		return nil
	}, funx.WithTimeout(time.Millisecond*50))
	assert.Equal(t, err, funx.ErrTimeout)
}

func TestWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Millisecond * 10)
		cancel()
	}()

	err := funx.DoWithTimeout(func() error {
		time.Sleep(time.Minute)
		return nil
	}, funx.WithTimeout(time.Second), funx.WithContext(ctx))
	assert.Equal(t, funx.ErrCanceled, err)
}
