package apns

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
)

// See "Create and Encrypt Your JSON Token" section in
// https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/establishing_a_token-based_connection_to_apns.

// Private key is not ECDSA key with P-256 curve.
var ErrJWTKeyNotECDSAP256 = errors.New("jwt: not ECDSA P-256 key")

// i2osp is an I2OSP (Integer to Octet Stream Primitive) function.
func i2osp(n *big.Int) []byte {
	return n.FillBytes(make([]byte, 32))
}

// es256 signs an input (JWS Signing Input)
// using Curve Digital Signature Algorithm (ECDSA)
// by key with P-256 curve type, specified IN RFC 7515.
// Returns digital signature (JWS Signature).
func es256(key *ecdsa.PrivateKey, input []byte) ([]byte, error) {
	if key.Curve != elliptic.P256() {
		return nil, ErrJWTKeyNotECDSAP256
	}

	h := sha256.New()
	_, err := h.Write(input)
	if err != nil {
		return nil, err
	}

	r, s, err := ecdsa.Sign(rand.Reader, key, h.Sum(nil))
	if err != nil {
		return nil, err
	}

	return append(i2osp(r), i2osp(s)...), nil
}

// GenerateBearer creates JSON token with keyID, teamID and issuedAt and
// encrypts by an authentication token signing key with the ES256 algorithm.
// Returns encrypted token that used in the authorization header of a notification request
// as bearer <token data>.
func GenerateBearer(key *ecdsa.PrivateKey, keyID, teamID string, issuedAt int64) (string, error) {
	// See RFC 7519 for JWT and RFC 7515 for JWS.
	header := fmt.Sprintf(`{"alg":"ES256","typ":"JWT","kid":"%s"}`, keyID) // JOSE Header (JWT Protected Header)
	payload := fmt.Sprintf(`{"iss":"%s","iat":%d}`, teamID, issuedAt)      // JWT Claims (JWS Payload)
	unsecured := base64.RawURLEncoding.EncodeToString([]byte(header)) + "." + base64.RawURLEncoding.EncodeToString([]byte(payload))
	sig, err := es256(key, []byte(unsecured)) // JWS Signature
	if err != nil {
		return "", err
	}
	t := unsecured + "." + base64.RawURLEncoding.EncodeToString([]byte(sig)) // JWS Compact Serialization
	return t, nil
}
