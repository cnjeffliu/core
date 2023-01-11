/*
 * @Author: cnzf1
 * @Date: 2021-12-15 14:18:20
 * @LastEditors: cnzf1
 * @LastEditTime: 2022-12-02 23:23:22
 * @Description:
 */
package thread

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	RunSafeEx(fn, nil)
}

// var err error
//
// RunSafeEx(func(){...}, &err)
func RunSafeEx(fn func(), err *error) {
	defer func() {
		res := Recover()
		if err != nil {
			if res != nil {
				*err = res
			} else {
				*err = nil
			}
		}
	}()

	fn()
}
