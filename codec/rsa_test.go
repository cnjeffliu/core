/*
 * @Author: Jeffrey.Liu
 * @Date: 2021-12-16 16:21:11
 * @LastEditors: Jeffrey.Liu
 * @LastEditTime: 2021-12-16 17:11:05
 * @Description:
 */
package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	priKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAujka1xH+mye8zIlAzFuQ6UsVLKWiNDsINBTDSPHZ3qQujakL
3LQQjYpiKMH5qE/inVJncZ/rQl1bspxSjncU4GLYu6GBcMwlbDNCpH1DZwXuZQlr
wCXghFnFVqdtnPJ28Oq55oDDjLQrQVdCGM79nwWnwpyUkkCl55JKinkWDHgZFRX4
6vRvu+pPEZm3q/LEdqLqFej2PgDiPrU5Xz/inD5JJpfWsALOQ2F2FrBx7IIK9Eat
Mfwp84p5wwqBaSi2l54CBPdbvKhL4K5cI+Mo3TKi7xtIWp+pfghVsFdZvP18I8K6
MXFALNc/GNRqZaHnu3glODVvHY0IRm/ON/nbZQIDAQABAoIBAQCjEJwDFdu3qx00
kT8vc0K6Niftd4BIciSlzkSOXFDmFyg4nX0onngcKL/5Zpmhm4oZLm4sXddYvn0s
MpxL6dRbA9M6wZqh1fEzBNPnS1S5IsV0rcIveDtYSW92iJeAJgSmwzNTtw8E50M1
LR5QsPf+xqn2zLuAMaHU3BHvnUYEVarFp53vtj0G4fXSTDNon8HEO3sIXvvKudqR
qe29ySIg+mS5qeBfF3JP63gDndSNT7/1JDhDPxWEm58S8FmwEaaU/RTOM0wxMUyV
jR/ZdYeJMir2JZLHC/5IrK67+nDTD0OuNqMWZzCkWEwEMU+SLJazMCkRX471wFoH
j8mhiaElAoGBAOc5fx4RkNS0hriXtsoWNCia0HnYQMu8vHb8xFC51lQdfA2slcX9
Y7AB6q6J3eJUi2l2fjzsUt9TANa6ab2O6g3KmWDT50Pxq0QFMOOPJVjDbn1QdGZA
A2AUp+YzPuWJ0/gf1FD0tTk5bbv3p7CB/muR4/TO6LPyDnK5mrSkZ9UXAoGBAM4t
Nr3ET3PEgHFGcn3KQzH14JU5A1Ma8sBP4RfKZyufU2I8ybN6RG4PcKHyzOB9/Gpb
607UsN4DdDQVSUdb8xaQXBZaQBGQMgqCAmjjqIQxy001jLvIHYcUxm26ZHyv0XwR
e3JkEwPqPaDa6LKlnhDGcs6G6J4B4aUR88EvmFjjAoGBAMySqyvwQKJgQh2JZRiw
wl72ceKLePCIwHnJsur1MHJlT79NZYmxYQR0/ayEn8JCKMIbKx89uyiI6GIStcEX
c27WRBNOB/uuEmfw68s5d8JrzhKjHwjkM9hLDi12Q3yUD+0kRBWIG9pQPA0k1MEu
kemcPwH2Gh4y16ObIQwXtSHrAoGAV+Hf5o2qDD+jPCV6IfI4KDCVNSYjK6Zd+OlT
mg91YJu+MC6XD0C7sGo2aWGUQNCS6kcaCvUQGuJAAv9bx+YCvQh1qDV5/8KGAgKe
wlTf/NE4xkVgIp7PL0gEuLrtoFRVJ9xP0Vek31NWR51n+NYthRsBztSkjM1igDkh
vKPr/V8CgYEAvSdqrkc1d5meeanRZjUPwcP2FGau2NHkaxcTlEYG6VWqpuGjoXEg
ANITv+uG8t+Uhkcx/xhborHEk/1xrqfpxeJQTPyGtwJHA3hR0/2NRA8RP5Cu0atZ
DOdjuNFpMQpZ4OQXYjYDU4GEgaCVBZkJEosz/kl32qQqidfyCS7IsbE=
-----END RSA PRIVATE KEY-----`
	pubKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAujka1xH+mye8zIlAzFuQ
6UsVLKWiNDsINBTDSPHZ3qQujakL3LQQjYpiKMH5qE/inVJncZ/rQl1bspxSjncU
4GLYu6GBcMwlbDNCpH1DZwXuZQlrwCXghFnFVqdtnPJ28Oq55oDDjLQrQVdCGM79
nwWnwpyUkkCl55JKinkWDHgZFRX46vRvu+pPEZm3q/LEdqLqFej2PgDiPrU5Xz/i
nD5JJpfWsALOQ2F2FrBx7IIK9EatMfwp84p5wwqBaSi2l54CBPdbvKhL4K5cI+Mo
3TKi7xtIWp+pfghVsFdZvP18I8K6MXFALNc/GNRqZaHnu3glODVvHY0IRm/ON/nb
ZQIDAQAB
-----END PUBLIC KEY-----`
	testBody = `this is the content`
)

func TestCryption(t *testing.T) {
	enc, err := NewRsaEncrypter([]byte(pubKey))
	assert.Nil(t, err)
	ret, err := enc.Encrypt([]byte(testBody))
	assert.Nil(t, err)

	dec, err := NewRsaDecrypter([]byte(priKey))
	assert.Nil(t, err)
	actual, err := dec.Decrypt(ret)
	assert.Nil(t, err)
	assert.Equal(t, testBody, string(actual))
}

func TestCryptionBase(t *testing.T) {
	enc, err := NewRsaEncrypter([]byte(pubKey))
	assert.Nil(t, err)
	ret, err := enc.EncryptBase64([]byte(testBody))
	assert.Nil(t, err)

	dec, err := NewRsaDecrypter([]byte(priKey))
	assert.Nil(t, err)
	actual, err := dec.DecryptBase64(string(ret))
	assert.Nil(t, err)
	assert.Equal(t, testBody, string(actual))
}

func TestBadPubKey(t *testing.T) {
	_, err := NewRsaEncrypter([]byte("foo"))
	assert.Equal(t, ErrPubKey, err)
}
