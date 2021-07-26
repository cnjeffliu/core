package algorithm

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// Decrypt
func Decrypt(privKey []byte, ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
