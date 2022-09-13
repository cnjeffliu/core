package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintFirst1BitNum(t *testing.T) {
	var s = uint64(8) // 1000
	d := PrintFisrt1BitNum(s)
	assert.Equal(t, d, uint64(8))

	s = 10 // 1010
	d = PrintFisrt1BitNum(s)
	assert.Equal(t, d, uint64(2))
}

func TestPrintFirst0BitNum(t *testing.T) {
	var s = uint64(8)
	d := PrintFisrt0BitNum(s)
	assert.Equal(t, d, uint64(1))

	s = 10
	d = PrintFisrt0BitNum(s)
	assert.Equal(t, d, uint64(1))
}
