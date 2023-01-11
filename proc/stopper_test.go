/*
 * @Author: cnzf1
 * @Date: 2022-12-02 21:02:42
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-01-11 11:08:38
 * @Description:
 */
package proc

import "testing"

func TestNopStopper(t *testing.T) {
	// no panic
	noopStopper.Stop()
}
