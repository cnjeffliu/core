package funx_test

import (
	"errors"
	"testing"

	"github.com/cnzf1/gocore/funx"
	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	assert.NotNil(t, funx.DoWithRetry(func() error {
		return errors.New("any")
	}))

	var times int
	assert.Nil(t, funx.DoWithRetry(func() error {
		times++
		if times == 3 {
			return nil
		}
		return errors.New("any")
	}))

	times = 0
	assert.NotNil(t, funx.DoWithRetry(func() error {
		times++
		if times == 3+1 {
			return nil
		}
		return errors.New("any")
	}))

	total := 2 * 3
	times = 0
	assert.Nil(t, funx.DoWithRetry(func() error {
		times++
		if times == total {
			return nil
		}
		return errors.New("any")
	}, funx.WithRetry(total)))
}
