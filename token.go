package apns

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
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

func AuthKeyFromFile(filename string) (*ecdsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return AuthKeyFromBytes(bytes)
}

func AuthKeyFromBytes(bytes []byte) (*ecdsa.PrivateKey, error) {
	b, _ := pem.Decode(bytes)
	if b == nil {
		return nil, errors.New("invalid .p8 PEM")
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
