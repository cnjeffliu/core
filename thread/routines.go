/*
 * @Author: Jeffrey Liu
 * @Date: 2021-12-15 14:18:20
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-08-27 16:06:04
 * @Description:
 */
package thread

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer Recover()

	fn()
}
