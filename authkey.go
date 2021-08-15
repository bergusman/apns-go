package apns

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

var (
	ErrAuthKeyBadPEM      = errors.New("authkey: invalid PEM")
	ErrAuthKeyBadPKCS8    = errors.New("authkey: invalid PKCS#8")
	ErrAuthKeyNotECDSA256 = errors.New("authkey: not ECDSA 256 bits")
)

// AuthKeyFromFile loads an authentication token signing key
// from the named text file with a .p8 file extension
// and returns an *ecdsa.PrivateKey.
func AuthKeyFromFile(name string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return AuthKeyFromBytes(data)
}

// AuthKeyFromBytes loads an authentication token signing key
// from the in memory bytes and returns an *ecdsa.PrivateKey.
func AuthKeyFromBytes(bytes []byte) (*ecdsa.PrivateKey, error) {
	p, _ := pem.Decode(bytes)
	if p == nil {
		return nil, ErrAuthKeyBadPEM
	}
	if len(p.Bytes) == 0 {
		return nil, ErrAuthKeyBadPEM
	}

	p8, err := x509.ParsePKCS8PrivateKey(p.Bytes)
	if err != nil {
		return nil, ErrAuthKeyBadPKCS8
	}

	key, ok := p8.(*ecdsa.PrivateKey)
	if !ok {
		return nil, ErrAuthKeyNotECDSA256
	}
	if key.Params().BitSize != 256 {
		return nil, ErrAuthKeyNotECDSA256
	}

	return key, nil
}
