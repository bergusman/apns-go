package apns

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"
)

func ExampleGenerateBearer() {

}

func TestI2OSP(t *testing.T) {
	n := big.NewInt(1337)
	h := hex.EncodeToString(i2osp(n, 32))
	if h != "0000000000000000000000000000000000000000000000000000000000000539" {
		t.Errorf("%q", h)
	}
}

func BenchmarkGenerateBearer(b *testing.B) {
	key, err := AuthKeyFromFile("testdata/AuthKeyValid.p8")
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
