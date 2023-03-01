/*
* @Author: cnzf1
* @Date: 2021-08-05 17:31:50
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-01 16:49:29
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
	"math/big"
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

// 生成RSA对象
func NewRsa(publicKey, privateKey string) (*Rsa, error) {
	rsaObj := &Rsa{
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	err := rsaObj.init()
	return rsaObj, err
}

// 生成pkcs1格式公钥私钥
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

// 生成pkcs8格式公钥私钥
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
			publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				return ErrPubKey
			}
		}
		this.rsaPublicKey = publicKey.(*rsa.PublicKey)
	}
	return nil
}

// 公钥加密
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

// 私钥解密
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

// 公钥加密+base64编码
func (this *Rsa) EncryptBase64(rawData []byte) (string, error) {
	enc, err := this.Encrypt(rawData)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(enc), nil
}

// 私钥解密+base64解码
func (this *Rsa) DecryptBase64(encData string) ([]byte, error) {
	dec, err := base64.StdEncoding.DecodeString(encData)
	if err != nil {
		return []byte(""), err
	}

	return this.Decrypt(dec)
}

// 私钥加密
func (this *Rsa) EncryptEx(rawData []byte) ([]byte, error) {
	return this.Sign(rawData, crypto.Hash(0))
}

// 公钥解密
func (this *Rsa) DecryptEx(encData []byte) ([]byte, error) {
	hashLen := 0
	var prefix []byte
	tLen := len(prefix) + hashLen
	k := (this.rsaPublicKey.N.BitLen() + 7) / 8
	if k < tLen+11 {
		return nil, ErrPubKeyNotRsa
	}

	c := new(big.Int).SetBytes(encData)
	m := encrypt(new(big.Int), this.rsaPublicKey, c)
	em := leftPad(m.Bytes(), k)
	out := unLeftPad(em)

	return out, nil
}

// 私钥加密+base64编码
func (this *Rsa) EncryptBase64Ex(rawData []byte) (string, error) {
	enc, err := this.EncryptEx(rawData)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(enc), nil
}

// 公钥解密+base64解码
func (this *Rsa) DecryptBase64Ex(encData string) ([]byte, error) {
	dec, err := base64.StdEncoding.DecodeString(encData)
	if err != nil {
		return []byte(""), err
	}

	return this.DecryptEx(dec)
}

// 私钥签名
func (this *Rsa) Sign(rawData []byte, algorithmSign crypto.Hash) ([]byte, error) {
	data := rawData
	if algorithmSign > 0 {
		hash := algorithmSign.New()
		hash.Write(rawData)
		data = hash.Sum(nil)
	}

	sign, err := rsa.SignPKCS1v15(rand.Reader, this.rsaPrivateKey, algorithmSign, data)
	if err != nil {
		return nil, err
	}
	return sign, err
}

// 公钥验签
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

// copy from crypt/rsa/pkcs1v15.go
func encrypt(c *big.Int, pub *rsa.PublicKey, m *big.Int) *big.Int {
	e := big.NewInt(int64(pub.E))
	c.Exp(m, e, pub.N)
	return c
}

// copy from crypt/rsa/pkcs1v15.go
func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out)-n:], input)
	return
}

func unLeftPad(input []byte) (out []byte) {
	n := len(input)
	t := 2
	for i := 2; i < n; i++ {
		if input[i] == 0xff {
			t = t + 1
		} else {
			if input[i] == input[0] {
				t = t + int(input[1])
			}
			break
		}
	}
	out = make([]byte, n-t)
	copy(out, input[t:])
	return
}
