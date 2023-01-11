/*
 * @Author: cnzf1
 * @Date: 2021-07-21 15:52:04
 * @Descripttion:
 */
package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"github.com/spaolacci/murmur3"
)

// Hash returns the hash value of data.
func Hash(data []byte) uint64 {
	return murmur3.Sum64(data)
}

// MD5
func MD5(input []byte) []byte {
	hash := md5.New()
	hash.Write(input)
	return hash.Sum(nil)
}

// SHA1
func SHA1(input []byte) []byte {
	hash := sha1.New()
	hash.Write(input)
	return hash.Sum(nil)
}

// SHA256
func SHA256(input []byte) []byte {
	hash := sha256.New()
	hash.Write(input)
	return hash.Sum(nil)
}

// SHA512
func SHA512(input []byte) []byte {
	hash := sha512.New()
	hash.Write(input)
	return hash.Sum(nil)
}
