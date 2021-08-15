package apns

import (
	"crypto/ecdsa"
	"errors"
	"net/http"
	"sync"
	"time"
)

// See https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/establishing_a_token-based_connection_to_apns.

// Token represents JSON Token used for token-based connection to APNs.
type Token struct {
	sync.Mutex

	// Authentication token signing key
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

func SetBearer(h http.Header, b string) {
	h.Set("authorization", "bearer "+b)
}
