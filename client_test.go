package apns

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientPush(t *testing.T) {
	key, err := AuthKeyFromFile("testdata/AuthKey_5MDQ4KLTY7.p8")
	if err != nil {
		t.Fatal(err)
	}
	token := NewToken(key, "5JZB9P77A7", "SUPERTEEM1")

	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		if r.Method != http.MethodPost {
			t.Errorf(":method: %v; want: %v", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/3/device/7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e" {
			t.Errorf(":path: %v, want: %v", r.URL.Path, "/3/device/7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e")
		}
		if r.Header.Get("authorization") != "bearer "+token.Bearer {
			t.Errorf("authorization: %v; want: %v", r.Header.Get("authorization"), "bearer "+token.Bearer)
		}
		if r.Header.Get("apns-id") != "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D" {
			t.Errorf("apns-id: %v; want: %v", r.Header.Get("apns-id"), "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D")
		}
		if r.Header.Get("apns-push-type") != "alert" {
			t.Errorf("apns-push-type: %v; want: %v", r.Header.Get("apns-push-type"), "alert")
		}
		if string(payload) != `{"aps":{"alert":{"title":"Hello"}}}` {
			t.Errorf("payload: %v; want: %v", string(payload), `{"aps":{"alert":{"title":"Hello"}}}`)
		}

		w.Header().Set("apns-id", "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D")
		w.WriteHeader(410)
		w.Write([]byte(`{"reason":"Unregistered","timestamp":1629000000}`))
	}))
	ts.EnableHTTP2 = true
	ts.StartTLS()
	defer ts.Close()

	client := NewClient(token, ts.Client())

	n := &Notification{
		DeviceToken: "7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e",
		Host:        ts.URL,
		ID:          "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D",
		PushType:    PushTypeAlert,
		Payload: BuildPayload(&APS{
			Alert: Alert{
				Title: "Hello",
			},
		}, nil),
	}

	res, err := client.Push(n)
	if err != nil {
		t.Fatal(err)
	}

	if res.ID != "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D" {
		t.Errorf("res.ID: %v; want: %v", res.ID, "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D")
	}
	if res.Status != Status410 {
		t.Errorf("res.Status: %v; want: %v", res.Status, Status410)
	}
	if res.Reason != ReasonUnregistered {
		t.Errorf("res.Reason: %v; want: %v", res.Reason, ReasonUnregistered)
	}
	if res.Timestamp != 1629000000 {
		t.Errorf("res.Timestamp: %v; want: %v", res.Timestamp, 1629000000)
	}
}

func TestClientPushErrors(t *testing.T) {
	client := NewClient(nil, nil)
	_, err := client.Push(nil)
	if err != nil {
		if err != ErrClientNotificationNil {
			t.Errorf("got: %v, want: ErrClientNotificationNil", err)
		}
	} else {
		t.Error("err must be not nil")
	}

	n := &Notification{}
	_, err = client.Push(n)
	if err != nil {
		if err != ErrClientTokenNil {
			t.Errorf("got: %v, want: ErrClientTokenNil", err)
		}
	} else {
		t.Error("err must be not nil")
	}
}
