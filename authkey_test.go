package apns

import (
	"os"
	"testing"
)

func TestAuthKeyNoSuchFile(t *testing.T) {
	_, err := AuthKeyFromFile("")
	if err.Error() != "open : no such file or directory" {
		t.Errorf("got: %q; want: %q", err, "open : no such file or directory")
	}
}

func TestAuthKeyBadPEM(t *testing.T) {
	names := []string{
		"testdata/AuthKeyBadPEM_1.p8",
		"testdata/AuthKeyBadPEM_2.p8",
	}
	for _, name := range names {
		_, err := AuthKeyFromFile(name)
		if err != nil {
			if err != ErrAuthKeyBadPEM {
				t.Errorf("%v: got: %q; want: ErrAuthKeyBadPEM", name, err)
			}
		} else {
			t.Errorf("%v: want not nil err", name)
		}
	}
}

func TestAuthKeyBadPKCS8(t *testing.T) {
	names := []string{
		"testdata/AuthKeyBadPKCS8_1.p8",
		"testdata/AuthKeyBadPKCS8_2.p8",
	}

	for _, name := range names {
		_, err := AuthKeyFromFile(name)
		if err != nil {
			if err != ErrAuthKeyBadPKCS8 {
				t.Errorf("%v: got: %q; want ErrAuthKeyBadPKCS8", name, err)
			}
		} else {
			t.Errorf("%v: want not nil err", name)
		}
	}
}

func TestAuthKeyNotECDSA256(t *testing.T) {
	names := []string{
		"testdata/AuthKeyRSA.p8",
		"testdata/AuthKeyED25519.p8",
		"testdata/AuthKeyECDSAP224.p8",
	}

	for _, name := range names {
		_, err := AuthKeyFromFile(name)
		if err != nil {
			if err != ErrAuthKeyNotECDSAP256 {
				t.Errorf("%v: got: %q; want ErrAuthKeyNotECDSAP256", name, err)
			}
		} else {
			t.Errorf("%v: want not nil err", name)
		}
	}
}

func TestAuthKey(t *testing.T) {
	key, err := AuthKeyFromFile("testdata/AuthKey_5MDQ4KLTY7.p8")
	if err != nil {
		t.Error(err)
	}
	if key == nil {
		t.Error("want not nil key")
	}
}

func BenchmarkAuthKeyFromBytes(b *testing.B) {
	file, err := os.ReadFile("testdata/AuthKey_5MDQ4KLTY7.p8")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AuthKeyFromBytes(file)
	}
}

func BenchmarkAuthKeyFromFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AuthKeyFromFile("testdata/AuthKey_5MDQ4KLTY7.p8")
	}
}
