package algorithm

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

// MD5
func MD5(input []byte) ([]byte, error) {
	hash := md5.New()
	_, err := hash.Write(input)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

// SHA1
func SHA1(input []byte) ([]byte, error) {
	hash := sha1.New()
	_, err := hash.Write(input)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

// SHA256
func SHA256(input []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write(input)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

// SHA512
func SHA512(input []byte) ([]byte, error) {
	hash := sha512.New()
	_, err := hash.Write(input)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
