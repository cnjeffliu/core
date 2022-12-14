/*
 * @Author: Jeffrey Liu
 * @Date: 2022-10-24 11:11:21
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-11-25 16:30:09
 * @Description:
 */
package mathx_test

import (
	"testing"

	"github.com/cnzf1/gocore/mathx"
	"github.com/stretchr/testify/assert"
)

func TestPrintFirst1BitNum(t *testing.T) {
	var s = uint64(8) // 1000
	d := mathx.PrintFisrt1BitNum(s)
	assert.Equal(t, d, uint64(8))

	s = 10 // 1010
	d = mathx.PrintFisrt1BitNum(s)
	assert.Equal(t, d, uint64(2))
}

func TestPrintFirst0BitNum(t *testing.T) {
	var s = uint64(8)
	d := mathx.PrintFisrt0BitNum(s)
	assert.Equal(t, d, uint64(1))

	s = 10
	d = mathx.PrintFisrt0BitNum(s)
	assert.Equal(t, d, uint64(1))
}
