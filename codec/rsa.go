/*
 * @Author: Jeffrey.Liu <zhifeng172@163.com>
 * @Date: 2021-08-05 17:31:50
 * @LastEditors: Jeffrey.Liu
 * @LastEditTime: 2021-12-16 17:15:59
 * @Description: 加密解密工具类
 */
package codec

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

var (
	ErrPrivKey   = errors.New("private key error")
	ErrPubKey    = errors.New("failed to parse PEM block containing the public key")
	ErrNotRsaKey = errors.New("key type is not RSA")
)

type (
	// RsaDecrypter represents a RSA decrypter.
	RsaDecrypter interface {
		Decrypt(input []byte) ([]byte, error)
		DecryptBase64(input string) ([]byte, error)
	}

	// RsaEncrypter represents a RSA encrypter.
	RsaEncrypter interface {
		Encrypt(input []byte) ([]byte, error)
		EncryptBase64(input []byte) (string, error)
	}

	rsaBase struct {
		bytesLimit int
	}

	rsaDecrypter struct {
		rsaBase
		privKey *rsa.PrivateKey
	}

	rsaEncrypter struct {
		rsaBase
		pubKey *rsa.PublicKey
	}
)

// NewRsaDecrypter returns a RsaDecrypter with the given privfile.
func NewRsaDecrypterWithFile(privfile string) (RsaDecrypter, error) {
	content, err := ioutil.ReadFile(privfile)
	if err != nil {
		return nil, err
	}

	return NewRsaDecrypter(content)
}

// NewRsaDecrypter returns a RsaDecrypter with the key content.
func NewRsaDecrypter(content []byte) (RsaDecrypter, error) {
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, ErrPrivKey
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &rsaDecrypter{
		rsaBase: rsaBase{
			bytesLimit: privKey.N.BitLen() >> 3,
		},
		privKey: privKey,
	}, nil
}

func (r *rsaDecrypter) Decrypt(input []byte) ([]byte, error) {
	return r.crypt(input, func(block []byte) ([]byte, error) {
		return rsaDecryptBlock(r.privKey, block)
	})
}

func (r *rsaDecrypter) DecryptBase64(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil
	}

	base64Decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, err
	}

	return r.Decrypt(base64Decoded)
}

func NewRsaEncrypterWithFile(pubfile string) (RsaEncrypter, error) {
	content, err := ioutil.ReadFile(pubfile)
	if err != nil {
		return nil, err
	}

	return NewRsaEncrypter(content)
}

// NewRsaEncrypter returns a RsaEncrypter with the given key.
func NewRsaEncrypter(key []byte) (RsaEncrypter, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrPubKey
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pubKey := pub.(type) {
	case *rsa.PublicKey:
		return &rsaEncrypter{
			rsaBase: rsaBase{
				// https://www.ietf.org/rfc/rfc2313.txt
				// The length of the data D shall not be more than k-11 octets, which is
				// positive since the length k of the modulus is at least 12 octets.
				bytesLimit: (pubKey.N.BitLen() >> 3) - 11,
			},
			pubKey: pubKey,
		}, nil
	default:
		return nil, ErrNotRsaKey
	}
}

func (r *rsaEncrypter) Encrypt(input []byte) ([]byte, error) {
	return r.crypt(input, func(block []byte) ([]byte, error) {
		return rsaEncryptBlock(r.pubKey, block)
	})
}

func (r *rsaEncrypter) EncryptBase64(input []byte) (string, error) {
	encryptedData, err := r.crypt(input, func(block []byte) ([]byte, error) {
		return rsaEncryptBlock(r.pubKey, block)
	})

	return base64.StdEncoding.EncodeToString(encryptedData), err
}

func (r *rsaBase) crypt(input []byte, cryptFn func([]byte) ([]byte, error)) ([]byte, error) {
	var result []byte
	inputLen := len(input)

	for i := 0; i*r.bytesLimit < inputLen; i++ {
		start := r.bytesLimit * i
		var stop int
		if r.bytesLimit*(i+1) > inputLen {
			stop = inputLen
		} else {
			stop = r.bytesLimit * (i + 1)
		}
		bs, err := cryptFn(input[start:stop])
		if err != nil {
			return nil, err
		}

		result = append(result, bs...)
	}

	return result, nil
}

func rsaDecryptBlock(privKey *rsa.PrivateKey, block []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privKey, block)
}

func rsaEncryptBlock(pubKey *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, msg)
}
