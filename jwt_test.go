package apns

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strings"
	"testing"
	"time"
)

func TestI2OSP(t *testing.T) {
	n := big.NewInt(1337)
	h := hex.EncodeToString(i2osp(n))
	if h != "0000000000000000000000000000000000000000000000000000000000000539" {
		t.Errorf("got: %v", h)
	}
}

func TestES256KeyNotP256(t *testing.T) {
	curves := []elliptic.Curve{elliptic.P224(), elliptic.P384(), elliptic.P521()}
	for _, curve := range curves {
		key, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		_, err = es256(key, []byte("input"))
		if err != nil {
			if err != ErrJWTKeyNotECDSAP256 {
				t.Errorf("curve %v: got: %q; want: ErrJWTKeyNot256Bits", curve.Params().Name, err)
			}
		} else {
			t.Errorf("curve %v: want not nil err", curve.Params().Name)
		}
	}
}

func TestES256Key(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	input := []byte("input")
	sig, err := es256(key, input)
	if err != nil {
		t.Fatal(err)
	}

	r, s := new(big.Int), new(big.Int)
	r.SetBytes(sig[:32])
	s.SetBytes(sig[32:])

	h := sha256.New()
	_, err = h.Write(input)
	if err != nil {
		t.Fatal(err)
	}

	if !ecdsa.Verify(&key.PublicKey, h.Sum(nil), r, s) {
		t.Error("cannot verify ES256 signature")
	}
}

func TestGenerateBearer(t *testing.T) {
	key, err := AuthKeyFromFile("testdata/AuthKey_5MDQ4KLTY7.p8")
	if err != nil {
		t.Fatal(err)
	}

	keyID := "5MDQ4KLTY7"
	teamID := "SUPERTEEM1"
	issuedAt := time.Now().Unix()

	token, err := GenerateBearer(key, keyID, teamID, issuedAt)
	if err != nil {
		t.Fatal(err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Error("token must have three parts spearated by dot")
	}

	unsecured := []byte(parts[0] + "." + parts[1])

	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		t.Fatal(err)
	}
	if len(sig) != 64 {
		t.Fatalf("len(sig): %v; want: 64", len(sig))
	}

	r, s := new(big.Int), new(big.Int)
	r.SetBytes(sig[:32])
	s.SetBytes(sig[32:])

	h := sha256.New()
	_, err = h.Write(unsecured)
	if err != nil {
		t.Fatal(err)
	}

	if !ecdsa.Verify(&key.PublicKey, h.Sum(nil), r, s) {
		t.Error("cannot verify ES256 signature")
	}

	// Header
	header64, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("cannot decode header from base64: %v", err)
	}

	var header struct {
		Alg string
		Typ string
		Kid string
	}
	err = json.Unmarshal(header64, &header)
	if err != nil {
		t.Fatalf("cannot decode header from json: %v", err)
	}

	if header.Alg != "ES256" {
		t.Errorf("header.alg: %v; want: %v", header.Alg, teamID)
	}
	if header.Kid != keyID {
		t.Errorf("header.kid: %v; want: %v", header.Kid, issuedAt)
	}

	// Payload
	payload64, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("cannot decode payload from base64: %v", err)
	}

	var payload struct {
		Iss string
		Iat int64
	}
	err = json.Unmarshal(payload64, &payload)
	if err != nil {
		t.Fatalf("cannot decode payload from json: %v", err)
	}

	if payload.Iss != teamID {
		t.Errorf("payload.iss: %v; want: %v", payload.Iss, teamID)
	}
	if payload.Iat != issuedAt {
		t.Errorf("payload.iat: %v; want: %v", payload.Iat, issuedAt)
	}
}

func BenchmarkGenerateBearer(b *testing.B) {
	key, err := AuthKeyFromFile("testdata/AuthKey_5MDQ4KLTY7.p8")
	if err != nil {
		b.Fatal(err)
	}

	keyID := "5MDQ4KLTY7"
	teamID := "SUPERTEEM1"
	issuedAt := time.Now().Unix()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateBearer(key, keyID, teamID, issuedAt)
	}
}
