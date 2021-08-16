package apns

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"net/http"
	"testing"
	"time"
)

func TestTokenKeyNil(t *testing.T) {
	token := NewToken(nil, "5JZB9P77A7", "SUPERTEEM1")
	_, err := token.Generate()
	if err != nil {
		if err != ErrTokenKeyNil {
			t.Errorf("err: %v; want: ErrTokenKeyNil", err)
		}
	} else {
		t.Error("err must be not nil")
	}
}

func TestTokenKeyInvalid(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	token := NewToken(key, "5JZB9P77A7", "SUPERTEEM1")
	_, err = token.Generate()
	if err != nil {
		if err != ErrJWTKeyNotECDSAP256 {
			t.Errorf("err: %v; want: ErrJWTKeyNotECDSAP256", err)
		}
	} else {
		t.Error("err must be not nil")
	}
}

func TestTokenExpired(t *testing.T) {
	token := &Token{}
	if !token.Expired() {
		t.Error("token must be expired")
	}

	token.IssuedAt = time.Now().Add(-(TokenRefreshInterval + 1) * time.Second).Unix()
	if !token.Expired() {
		t.Error("token must be expired")
	}

	token.RefreshInterval = 60
	token.IssuedAt = time.Now().Add(-61 * time.Second).Unix()
	if !token.Expired() {
		t.Error("token must be expired")
	}
}

func TestTokenNotExpired(t *testing.T) {
	token := &Token{
		IssuedAt: time.Now().Unix(),
	}
	if token.Expired() {
		t.Error("token must be not expired")
	}

	token.IssuedAt = time.Now().Add(-(TokenRefreshInterval - 1) * time.Second).Unix()
	if token.Expired() {
		t.Error("token must be not expired")
	}

	token.RefreshInterval = 3600
	token.IssuedAt = time.Now().Add(-3599 * time.Second).Unix()
	if token.Expired() {
		t.Error("token must be not expired")
	}
}

func TestTokenGenerateIfExpired(t *testing.T) {
	key, err := AuthKeyFromFile("testdata/AuthKey_5MDQ4KLTY7.p8")
	if err != nil {
		t.Fatal(err)
	}

	token := NewToken(key, "5JZB9P77A7", "SUPERTEEM1")

	want, err := token.GenerateIfExpired()
	if err != nil {
		t.Fatal(err)
	}

	got, err := token.GenerateIfExpired()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got: %v; want: %v", got, want)
	}
}

func TestSetAuthorization(t *testing.T) {
	key, err := AuthKeyFromFile("testdata/AuthKey_5MDQ4KLTY7.p8")
	if err != nil {
		t.Fatal(err)
	}

	token := NewToken(key, "5JZB9P77A7", "SUPERTEEM1")
	h := make(http.Header)
	token.SetAuthorization(h)

	if h.Get("authorization") != "bearer "+token.Bearer {
		t.Errorf("invalid authorization header: %v", h.Get("authorization"))
	}
}

func TestSetAuthorizationWithInvalidKey(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	token := NewToken(key, "5JZB9P77A7", "SUPERTEEM1")
	h := make(http.Header)
	err = token.SetAuthorization(h)
	if err != nil {
		if err != ErrJWTKeyNotECDSAP256 {
			t.Errorf("err: %v; want: ErrJWTKeyNotECDSAP256", err)
		}
	} else {
		t.Error("err must be not nil")
	}
}

func TestSetBearer(t *testing.T) {
	bearer := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjVKWkI5UDc3QTcifQ.eyJpc3MiOiJTVVBFUlRFRU0xIiwiaWF0IjoxNjI5MDAwMDAwfQ.f9paE35q8NzQyrnJaDbYb1yBzBcj8_jHpeMe06y5OERf-0ERD70VZCS6dG-toWKSxeevYEARKMwoGliJtZ7atQ"
	h := make(http.Header)
	SetBearer(h, bearer)
	if h.Get("authorization") != "bearer "+bearer {
		t.Errorf("invalid authorization header: %v", h.Get("authorization"))
	}
}
