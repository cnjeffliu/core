/*
 * @Author: cnzf1
 * @Date: 2022-11-29 13:54:49
 * @LastEditors: cnzf1
 * @LastEditTime: 2022-12-02 18:19:58
 * @Description:
 */
package proc_test

import (
	"testing"

	"github.com/cnzf1/gocore/proc"
	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {
	profiler := proc.StartProfile()
	// start again should not work
	assert.NotNil(t, proc.StartProfile())
	profiler.Stop()
	// stop twice
	profiler.Stop()
}
