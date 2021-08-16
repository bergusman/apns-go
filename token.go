package apns

import (
	"crypto/ecdsa"
	"errors"
	"net/http"
	"sync"
	"time"
)

// See https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/establishing_a_token-based_connection_to_apns.

var ErrTokenKeyNil = errors.New("token: key is nil")

// For security, APNs requires you to refresh your token regularly.
// Refresh your token no more than once every 20 minutes
// and no less than once every 60 minutes.
// APNs rejects any request whose token
// contains a timestamp that is more than one hour old.
// Similarly, APNs reports an error if you recreate your tokens
// more than once every 20 minutes.
const TokenRefreshInterval = 2400 // 40 minutes

// Token represents JSON Token used for token-based connection to APNs.
type Token struct {
	sync.Mutex

	// Authentication token signing key extracted from a text file (with a .p8 file extension).
	// Use AuthKeyFromFile or AuthKeyFromBytes.
	Key *ecdsa.PrivateKey

	// The 10-character Key ID you obtained from your developer account
	// with an authentication token signing key.
	KeyID string

	// The value is the 10-character Team ID you use for developing your company’s apps.
	// Obtain this value from your developer account.
	TeamID string

	// Time window in from IssuedAt to IssuedAt + RefreshInterval
	// when Bearer is not expired.
	RefreshInterval int64 // in seconds

	// Time at Bearer was generated.
	IssuedAt int64 // Epoch time in seconds

	// Generated JWT Token for APNs request authorization at IssuedAt time.
	Bearer string
}

// NewToken returns Token with key, keyID, teamID with default TokenRefreshInterval.
// How obtain credentials,
// see https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/establishing_a_token-based_connection_to_apns#2943371
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
		return time.Now().Unix() > t.IssuedAt+TokenRefreshInterval
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
		return "", ErrTokenKeyNil
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
