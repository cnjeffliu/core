package algorithm

import (
	"testing"
)

func TestHash(t *testing.T) {
	input := "abcdefghijklmnopqrstuvwxyz"
	var result []byte
	var err error

	//MD5算法.
	for i := 0; i < 5000; i++ {
		result, err = MD5([]byte(input))
		if err != nil {
			t.Log(err)
			return
		}
	}
	t.Logf("md5: %x\n", result)

	//SHA1算法.
	result, err = SHA1([]byte(input))
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("sha1: %x\n", result)

	//SHA256算法.
	result, err = SHA256([]byte(input))
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("sha256: %x\n", result)

	//SHA512算法.
	result, err = SHA512([]byte(input))
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("sha512: %x\n", result)

}
