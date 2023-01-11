/*
 * @Author: cnzf1
 * @Date: 2022-12-02 21:02:42
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-01-11 11:08:51
 * @Description:
 */
package proc

var noopStopper nilStopper

type (
	// Stopper interface wraps the method Stop.
	Stopper interface {
		Stop()
	}

	nilStopper struct{}
)

func (ns nilStopper) Stop() {
}
