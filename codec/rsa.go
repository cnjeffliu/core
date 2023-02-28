/*
* @Author: cnzf1
* @Date: 2021-08-05 17:31:50
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-02-28 19:54:37
* @Description:
*/
package codec

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"

	"github.com/cnzf1/gocore/lang"
)

var (
	ErrPrivKey       = errors.New("private key error")
	ErrPubKey        = errors.New("failed to parse PEM block containing the public key")
	ErrPrivKeyNotRsa = errors.New("private key type is not RSA")
	ErrPubKeyNotRsa  = errors.New("public key type is not RSA")
)

type Rsa struct {
	privateKey    string
	publicKey     string
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
}

func NewRsa(publicKey, privateKey string) (*Rsa, error) {
	rsaObj := &Rsa{
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	err := rsaObj.init()
	return rsaObj, err
}

/**
 * 生成pkcs1格式公钥私钥
 */
func CreateRsa(keyLength int) (*Rsa, error) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return nil, err
	}

	privateKey := string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(rsaPrivateKey),
	}))

	derPkix, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	publicKey := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}))

	return NewRsa(publicKey, privateKey)
}

/**
 * 生成pkcs8格式公钥私钥
 */
func CreateRsaPkcs8(keyLength int) (*Rsa, error) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return nil, err
	}

	privateKey := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: MarshalPKCS8PrivateKey(rsaPrivateKey),
	}))

	derPkix, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	publicKey := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}))

	return NewRsa(publicKey, privateKey)
}

func (this *Rsa) init() (err error) {
	var privateKey lang.AnyType
	var publicKey lang.AnyType

	if this.privateKey != "" {
		block, _ := pem.Decode([]byte(this.privateKey))
		if block == nil || strings.Index(block.Type, "PRIVATE KEY") < 0 {
			return ErrPrivKeyNotRsa
		}

		//pkcs1
		if strings.Index(this.privateKey, "BEGIN RSA") > 0 {
			this.rsaPrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return ErrPrivKey
			}
		} else { //pkcs8
			privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return ErrPrivKey
			}
			this.rsaPrivateKey = privateKey.(*rsa.PrivateKey)
		}
	}

	if this.publicKey != "" {
		block, _ := pem.Decode([]byte(this.publicKey))
		if block == nil || strings.Index(block.Type, "PUBLIC KEY") < 0 {
			return ErrPubKeyNotRsa
		}

		publicKey, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return ErrPubKey
		}
		this.rsaPublicKey = publicKey.(*rsa.PublicKey)
	}
	return nil
}

/**
 * 加密
 */
func (this *Rsa) Encrypt(rawData []byte) ([]byte, error) {
	blockLength := this.rsaPublicKey.N.BitLen()/8 - 11
	if len(rawData) <= blockLength {
		return rsa.EncryptPKCS1v15(rand.Reader, this.rsaPublicKey, []byte(rawData))
	}

	buffer := bytes.NewBufferString("")

	pages := len(rawData) / blockLength

	for index := 0; index <= pages; index++ {
		start := index * blockLength
		end := (index + 1) * blockLength
		if index == pages {
			if start == len(rawData) {
				continue
			}
			end = len(rawData)
		}

		chunk, err := rsa.EncryptPKCS1v15(rand.Reader, this.rsaPublicKey, rawData[start:end])
		if err != nil {
			return nil, err
		}
		buffer.Write(chunk)
	}
	return buffer.Bytes(), nil
}

/**
 * 加密
 */
func (this *Rsa) EncryptBase64(rawData []byte) (string, error) {
	enc, err := this.Encrypt(rawData)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(enc), nil
}

/**
 * 解密
 */
func (this *Rsa) Decrypt(encData []byte) ([]byte, error) {
	blockLength := this.rsaPublicKey.N.BitLen() / 8
	if len(encData) <= blockLength {
		return rsa.DecryptPKCS1v15(rand.Reader, this.rsaPrivateKey, encData)
	}

	buffer := bytes.NewBufferString("")

	pages := len(encData) / blockLength
	for index := 0; index <= pages; index++ {
		start := index * blockLength
		end := (index + 1) * blockLength
		if index == pages {
			if start == len(encData) {
				continue
			}
			end = len(encData)
		}

		chunk, err := rsa.DecryptPKCS1v15(rand.Reader, this.rsaPrivateKey, encData[start:end])
		if err != nil {
			return nil, err
		}
		buffer.Write(chunk)
	}
	return buffer.Bytes(), nil
}

/**
 * 解密
 */
func (this *Rsa) DecryptBase64(encData string) ([]byte, error) {
	dec, err := base64.StdEncoding.DecodeString(encData)
	if err != nil {
		return []byte(""), err
	}

	return this.Decrypt(dec)
}

/**
 * 签名
 */
func (this *Rsa) Sign(rawData []byte, algorithmSign crypto.Hash) ([]byte, error) {
	hash := algorithmSign.New()
	hash.Write(rawData)
	sign, err := rsa.SignPKCS1v15(rand.Reader, this.rsaPrivateKey, algorithmSign, hash.Sum(nil))
	if err != nil {
		return nil, err
	}
	return sign, err
}

/**
 * 验签
 */
func (this *Rsa) Verify(rawData []byte, sign []byte, algorithmSign crypto.Hash) bool {
	h := algorithmSign.New()
	h.Write(rawData)
	return rsa.VerifyPKCS1v15(this.rsaPublicKey, algorithmSign, h.Sum(nil), sign) == nil
}

func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)
	k, _ := asn1.Marshal(info)
	return k
}
