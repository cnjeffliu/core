/*
 * @Author: Jeffrey.Liu
 * @Date: 2021-10-15 17:30:22
 * @Last Modified by: Jeffrey.Liu
 * @Last Modified time: 2021-10-15 17:33:25
 */
package math

// 从低位开始，获取s中的第一个为1的位所代表的整数
// 比如s=3（0011）返回1; 若s=6(0110)返回2
func PrintFisrt1BitNum(s uint64) (d uint64) {
	return s & (^(s - 1))
}

// 从低位开始，获取s中的第一个为0的位所代表的整数
// 比如s=3（0011）返回4; 若s=6(0110)返回1
func PrintFisrt0BitNum(s uint64) (d uint64) {
	s = ^s
	return PrintFisrt1BitNum(s)
}
