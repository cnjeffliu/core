package codec

import "crypto"

type (
	// Decrypter represents a RSA decrypter.
	Decrypter interface {
		Decrypt(encData []byte) ([]byte, error)
		DecryptBase64(encData string) ([]byte, error)
	}

	// Encrypter represents a RSA encrypter.
	Encrypter interface {
		Encrypt(rawData []byte) ([]byte, error)
		EncryptBase64(rawData []byte) (string, error)
	}

	Crypter interface {
		Encrypter
		Decrypter
	}

	Signer interface {
		Sign(rawData []byte, algorithmSign crypto.Hash) ([]byte, error)
		Verify(rawData []byte, sign []byte, algorithmSign crypto.Hash) bool
	}
)
