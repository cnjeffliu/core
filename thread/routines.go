/*
 * @Author: Jeffrey Liu
 * @Date: 2021-12-15 14:18:20
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-12-02 22:13:04
 * @Description:
 */
package thread

func GoSafe(fn func()) {
	go RunSafe(fn)
}

// var err error
//
// GoSafeEx(func(){...}, &err)
func GoSafeEx(fn func(), err *error) {
	go RunSafeEx(fn, nil)
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
