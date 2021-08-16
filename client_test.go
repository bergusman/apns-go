package apns

import (
	"testing"
)

func TestClientPush(t *testing.T) {
	n := &Notification{
		Host: "https://example.com",
	}

	key, err := AuthKeyFromFile("testdata/AuthKey_5MDQ4KLTY7.p8")
	if err != nil {
		t.Fatal(err)
	}

	token := NewToken(key, "5JZB9P77A7", "SUPERTEEM1")

	client := NewClient(token, nil)
	client.Push(n)
}
