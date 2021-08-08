package apns

import (
	"testing"
)

func TestAuthKeyValid(t *testing.T) {
	key, err := AuthKeyFromFile("testdata/AuthKeyValid.p8")
	if err != nil {
		t.Error(err)
	}
	if key == nil {
		t.Error("Key must not be nil")
	}
}

func TestAuthKeyNoSuchFile(t *testing.T) {
	_, err := AuthKeyFromFile("")
	if err.Error() != "open : no such file or directory" {
		t.Error(`Want "no such file" error`)
	}
}

func TestAuthKeyEmpty(t *testing.T) {
	_, err := AuthKeyFromFile("testdata/AuthKeyEmpty.p8")
	if err != nil {
		if err.Error() != ErrAuthKeyBadPEM.Error() {
			t.Errorf("err = %q; want ErrAuthKeyBadPEM", err.Error())
		}
	} else {
		t.Error("Want err not nil")
	}
}
