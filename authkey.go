package apns

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// Authentication Token Signing Key
// See "Obtain an Encryption Key and Key ID from Apple" secion in
// https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/establishing_a_token-based_connection_to_apns.

var (
	// .p8 file is not PEM encoded or invalid PEM file.
	ErrAuthKeyBadPEM = errors.New("authkey: invalid PEM")

	// PEM block does not contain private key in PKCS #8, ASN.1 DER form.
	ErrAuthKeyBadPKCS8 = errors.New("authkey: invalid PKCS#8")

	// Authentication token signing key is not ECDSA private key with P-256 curve.
	ErrAuthKeyNotECDSAP256 = errors.New("authkey: not ECDSA P-256")
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
		return nil, ErrAuthKeyNotECDSAP256
	}
	if key.Curve != elliptic.P256() {
		return nil, ErrAuthKeyNotECDSAP256
	}

	return key, nil
}
