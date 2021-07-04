package apns

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"
)

func i2osp(n *big.Int, l int) []byte {
	r := make([]byte, l)
	b := n.Bytes()
	copy(r[l-len(b):], b)
	return r
}

func es256(key *ecdsa.PrivateKey, body string) ([]byte, error) {
	h := crypto.SHA256.New()
	_, err := h.Write([]byte(body))
	if err != nil {
		return nil, err
	}

	r, s, err := ecdsa.Sign(rand.Reader, key, h.Sum(nil))
	if err != nil {
		return nil, err
	}

	return append(i2osp(r, 32), i2osp(s, 32)...), nil
}

func GenerateBearer(key *ecdsa.PrivateKey, keyId, teamId string) (string, int64, error) {
	iat := time.Now().Unix()
	h := fmt.Sprintf(`{"alg":"ES256","typ":"JWT","kid":"%s"}`, keyId)
	p := fmt.Sprintf(`{"iss":"%s","iat":%d}`, teamId, iat)
	hp := base64.RawURLEncoding.EncodeToString([]byte(h)) + "." + base64.RawURLEncoding.EncodeToString([]byte(p))
	sig, err := es256(key, hp)
	if err != nil {
		return "", 0, err
	}
	t := hp + "." + base64.RawURLEncoding.EncodeToString([]byte(sig))
	return t, iat, nil
}
