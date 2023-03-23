package mapx_test

import (
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/cnzf1/gocore/collection/mapx"
	"github.com/stretchr/testify/assert"
)

func TestSafeMap(t *testing.T) {
	tests := []struct {
		size      int
		exception int
	}{
		{
			100000,
			2000,
		},
		{
			100000,
			50,
		},
	}
	for _, test := range tests {
		t.Run(strconv.Itoa(test.size), func(t *testing.T) {
			testSafeMapWithParameters(t, test.size, test.exception)
		})
	}
}

func TestSafeMap_CopyNew(t *testing.T) {
	const (
		size       = 100000
		exception1 = 5
		exception2 = 500
	)
	m := mapx.NewSafeMap()

	for i := 0; i < size; i++ {
		m.Set(i, i)
	}
	for i := 0; i < size; i++ {
		if i%exception1 == 0 {
			m.Remove(i)
		}
	}

	for i := size; i < size<<1; i++ {
		m.Set(i, i)
	}
	for i := size; i < size<<1; i++ {
		if i%exception2 != 0 {
			m.Remove(i)
		}
	}

	for i := 0; i < size; i++ {
		val, ok := m.Get(i)
		if i%exception1 != 0 {
			assert.True(t, ok)
			assert.Equal(t, i, val.(int))
		} else {
			assert.False(t, ok)
		}
	}
	for i := size; i < size<<1; i++ {
		val, ok := m.Get(i)
		if i%exception2 == 0 {
			assert.True(t, ok)
			assert.Equal(t, i, val.(int))
		} else {
			assert.False(t, ok)
		}
	}
}

func testSafeMapWithParameters(t *testing.T, size, exception int) {
	m := mapx.NewSafeMap()

	for i := 0; i < size; i++ {
		m.Set(i, i)
	}
	for i := 0; i < size; i++ {
		if i%exception != 0 {
			m.Remove(i)
		}
	}

	assert.Equal(t, size/exception, m.Size())

	for i := size; i < size<<1; i++ {
		m.Set(i, i)
	}
	for i := size; i < size<<1; i++ {
		if i%exception != 0 {
			m.Remove(i)
		}
	}

	for i := 0; i < size<<1; i++ {
		val, ok := m.Get(i)
		if i%exception == 0 {
			assert.True(t, ok)
			assert.Equal(t, i, val.(int))
		} else {
			assert.False(t, ok)
		}
	}
}

func TestSafeMap_Range(t *testing.T) {
	const (
		size       = 100000
		exception1 = 5
		exception2 = 500
	)

	m := mapx.NewSafeMap()
	newMap := mapx.NewSafeMap()

	for i := 0; i < size; i++ {
		m.Set(i, i)
	}
	for i := 0; i < size; i++ {
		if i%exception1 == 0 {
			m.Remove(i)
		}
	}

	for i := size; i < size<<1; i++ {
		m.Set(i, i)
	}
	for i := size; i < size<<1; i++ {
		if i%exception2 != 0 {
			m.Remove(i)
		}
	}

	var count int32
	m.Range(func(k, v any) bool {
		atomic.AddInt32(&count, 1)
		newMap.Set(k, v)
		return true
	})
	assert.Equal(t, int(atomic.LoadInt32(&count)), m.Size())
	// assert.Equal(t, m.dirtyNew, newMap.dirtyNew)
	// assert.Equal(t, m.dirtyOld, newMap.dirtyOld)
}
