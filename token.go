package apns

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	ErrAuthKeyBadPEM = errors.New("invalid PEM")
)

type Token struct {
	sync.Mutex
	Key             *ecdsa.PrivateKey
	KeyID           string
	TeamID          string
	RefreshInterval int64
	IssuedAt        int64
	Bearer          string
}

func NewToken(key *ecdsa.PrivateKey, keyID, teamID string) *Token {
	return &Token{
		Key:    key,
		KeyID:  keyID,
		TeamID: teamID,
	}
}

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

// AuthKeyFromBytes load an authentication token signing key
// from the in memory bytes and returns an *ecdsa.PrivateKey.
func AuthKeyFromBytes(bytes []byte) (*ecdsa.PrivateKey, error) {
	b, _ := pem.Decode(bytes)
	if b == nil {
		return nil, ErrAuthKeyBadPEM
	}

	p8, err := x509.ParsePKCS8PrivateKey(b.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := p8.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("not ECDSA private key")
	}
	return key, nil
}

func (t *Token) Expired() bool {
	if t.RefreshInterval > 0 {
		return time.Now().Unix() > t.IssuedAt+t.RefreshInterval
	} else {
		return time.Now().Unix() > t.IssuedAt+2400 // 40 minutes
	}
}

func (t *Token) GenerateIfExpired() (string, error) {
	t.Lock()
	defer t.Unlock()
	if t.Expired() {
		return t.Generate()
	}
	return t.Bearer, nil
}

func (t *Token) Generate() (string, error) {
	if t.Key == nil {
		return "", errors.New("key is nil")
	}
	issuedAt := time.Now().Unix()
	bearer, err := GenerateBearer(t.Key, t.KeyID, t.TeamID, issuedAt)
	if err != nil {
		return "", err
	}
	t.Bearer = bearer
	t.IssuedAt = issuedAt
	return bearer, nil
}

func (t *Token) SetAuthorization(h http.Header) error {
	bearer, err := t.GenerateIfExpired()
	if err != nil {
		return err
	}
	h.Set("authorization", "bearer "+bearer)
	return nil
}
