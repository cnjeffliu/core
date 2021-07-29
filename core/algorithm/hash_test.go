/*
 * @Author: Jeffrey <zhifeng172@163.com>
 * @Date: 2021-07-21 15:55:09
 * @Descripttion:
 */
package algorithm

import (
	"testing"
)

func TestHash(t *testing.T) {
	input := "abcdefghijklmnopqrstuvwxyz"
	var result []byte

	//MD5算法.
	result = MD5([]byte(input))
	t.Logf("md5: %x\n", result)

	//SHA1算法.
	result = SHA1([]byte(input))
	t.Logf("sha1: %x\n", result)

	//SHA256算法.
	result = SHA256([]byte(input))
	t.Logf("sha256: %x\n", result)

	//SHA512算法.
	result = SHA512([]byte(input))
	t.Logf("sha512: %x\n", result)

}
