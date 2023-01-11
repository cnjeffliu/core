/*
 * @Author: cnzf1
 * @Date: 2021-07-21 15:55:09
 * @Descripttion:
 */
package hash

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	input := "abcdefghijklmnopqrstuvwxyz"

	result := Hash([]byte(input))
	assert.Equal(t, "749c9d7e516f4aa9", fmt.Sprintf("%x", result))
}

func TestMD5(t *testing.T) {
	input := "abcdefghijklmnopqrstuvwxyz"
	var result []byte

	//MD5算法.
	result = MD5([]byte(input))
	assert.Equal(t, "c3fcd3d76192e4007dfb496cca67e13b", fmt.Sprintf("%x", result))
}

func TestSHA1(t *testing.T) {
	input := "abcdefghijklmnopqrstuvwxyz"
	//SHA1算法.
	result := SHA1([]byte(input))
	t.Logf("sha1: %x\n", result)

}

func TestSHA256(t *testing.T) {
	input := "abcdefghijklmnopqrstuvwxyz"

	//SHA256算法.
	result := SHA256([]byte(input))
	assert.Equal(t, "71c480df93d6ae2f1efad1447c66c9525e316218cf51fc8d9ed832f2daf18b73", fmt.Sprintf("%x", result))

}

func TestSHA512(t *testing.T) {
	input := "abcdefghijklmnopqrstuvwxyz"

	//SHA512算法.
	result := SHA512([]byte(input))
	assert.Equal(t, "4dbff86cc2ca1bae1e16468a05cb9881c97f1753bce3619034898faa1aabe429955a1bf8ec483d7421fe3c1646613a59ed5441fb0f321389f77f48a879c7b1f1", fmt.Sprintf("%x", result))

}
