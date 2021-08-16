package apns

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestNotificationPayload(t *testing.T) {
	tests := []struct {
		payload interface{}
		want    []byte
	}{
		{
			payload: `{"aps":{"alert":{"title":"Hello"}}}`,
			want:    []byte(`{"aps":{"alert":{"title":"Hello"}}}`),
		},
		{
			payload: []byte(`{"aps":{"alert":{"title":"Hello"}}}`),
			want:    []byte(`{"aps":{"alert":{"title":"Hello"}}}`),
		},
		{
			payload: map[string]interface{}{
				"aps": map[string]interface{}{
					"alert": map[string]interface{}{
						"title": "Hello",
					},
				},
			},
			want: []byte(`{"aps":{"alert":{"title":"Hello"}}}`),
		},
		{
			payload: BuildPayload(&APS{
				Alert: Alert{
					Title: "Hello",
				},
			}, nil),
			want: []byte(`{"aps":{"alert":{"title":"Hello"}}}`),
		},
	}

	for _, tt := range tests {
		n := &Notification{
			Payload: tt.payload,
		}

		b, err := json.Marshal(n)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != string(tt.want) {
			t.Errorf("got: %v; want: %v", string(b), string(tt.want))
		}
	}
}

func TestNotificationPath(t *testing.T) {
	n := &Notification{}
	if n.Path() != "/3/device/" {
		t.Errorf("got: %v; want: %v", n.Path(), "/3/device/")
	}

	n = &Notification{
		DeviceToken: "7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e",
	}
	if n.Path() != "/3/device/7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e" {
		t.Errorf("got: %v; want: %v", n.Path(), "/3/device/7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e")
	}
}

func TestNotificationURL(t *testing.T) {
	n := &Notification{}
	if n.URL() != "https://api.push.apple.com/3/device/" {
		t.Errorf("got: %v; want: %v", n.URL(), "https://api.push.apple.com/3/device/")
	}

	n = &Notification{
		Host: HostDevelopment,
	}
	if n.URL() != "https://api.sandbox.push.apple.com/3/device/" {
		t.Errorf("got: %v; want: %v", n.URL(), "https://api.sandbox.push.apple.com/3/device/")
	}

	n = &Notification{
		Host:        HostProductionPort2197,
		DeviceToken: "7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e",
	}
	if n.URL() != "https://api.push.apple.com:2197/3/device/7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e" {
		t.Errorf("got: %v; want: %v", n.URL(), "https://api.push.apple.com:2197/3/device/7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e")
	}
}

func TestNotificationSetHeaders(t *testing.T) {
	n := &Notification{}
	h := make(http.Header)
	n.SetHeaders(h)
	if len(h) != 0 {
		t.Error("header must be empty")
	}

	n = &Notification{
		ID:         "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D",
		Topic:      "com.example.app",
		PushType:   PushTypeAlert,
		Expiration: "1629000000",
		Priority:   10,
		CollapseID: "448C1AA7-B421-44D2-A995-2E4A7F1AE29E",
	}
	h = make(http.Header)
	n.SetHeaders(h)

	if h.Get("apns-id") != "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D" {
		t.Errorf("apns-id: %v; want: %v", h.Get("apns-id"), "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D")
	}
	if h.Get("apns-topic") != "com.example.app" {
		t.Errorf("apns-topic: %v; want: %v", h.Get("apns-topic"), "com.example.app")
	}
	if h.Get("apns-push-type") != "alert" {
		t.Errorf("apns-push-type: %v; want: %v", h.Get("apns-push-type"), "alert")
	}
	if h.Get("apns-expiration") != "1629000000" {
		t.Errorf("apns-expiration: %v; want: %v", h.Get("apns-expiration"), "1629000000")
	}
	if h.Get("apns-priority") != "10" {
		t.Errorf("apns-priority: %v; want: %v", h.Get("apns-priority"), "10")
	}
	if h.Get("apns-collapse-id") != "448C1AA7-B421-44D2-A995-2E4A7F1AE29E" {
		t.Errorf("apns-collapse-id: %v; want: %v", h.Get("apns-collapse-id"), "448C1AA7-B421-44D2-A995-2E4A7F1AE29E")
	}
}

func TestNotificationBuildRequest(t *testing.T) {
	n := &Notification{
		DeviceToken: "7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e",
		ID:          "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D",
		Topic:       "com.example.app",
		Payload:     `{"aps":{"alert":{"title":"Hello"}}}`,
	}
	req, err := n.BuildRequest()
	if err != nil {
		t.Fatal(err)
	}

	if req.Method != http.MethodPost {
		t.Errorf("req.Method: %v; want: %v", req.Method, http.MethodPost)
	}
	if req.URL.String() != "https://api.push.apple.com/3/device/7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e" {
		t.Errorf("req.URL: %v; want: %v", req.URL, "https://api.push.apple.com/3/device/7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e")
	}
	if req.Header.Get("apns-id") != "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D" {
		t.Errorf("req.Header[apns-id]: %v; want: %v", req.Header.Get("apns-id"), "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D")
	}
	if req.Header.Get("apns-topic") != "com.example.app" {
		t.Errorf("req.Header[apns-topic]: %v; want: %v", req.Header.Get("apns-topic"), "com.example.app")
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != `{"aps":{"alert":{"title":"Hello"}}}` {
		t.Errorf("req.Body: %v; want: %v", string(body), `{"aps":{"alert":{"title":"Hello"}}}`)
	}
}

func TestNotificationBuildRequestErrors(t *testing.T) {
	n := &Notification{
		Host: ":::::",
	}
	_, err := n.BuildRequest()
	if err == nil {
		t.Error("err must be not nil")
	}

	n = &Notification{
		Payload: "{{{",
	}
	_, err = n.BuildRequest()
	if err == nil {
		t.Error("err must be not nil")
	}
}

func BenchmarkNotificationSimplePayload(b *testing.B) {
	n := &Notification{
		Payload: BuildPayload(&APS{
			Alert: Alert{
				Title: "Title",
				Body:  "Body",
			},
		}, nil),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(n)
	}
}

func BenchmarkNotificationFullPayload(b *testing.B) {
	n := &Notification{
		Payload: BuildPayload(&APS{
			Alert: Alert{
				Title:           "Title",
				TitleLocKey:     "TitleLocKey",
				TitleLocArgs:    []string{"Arg1", "Arg2"},
				Subtitle:        "Subtitle",
				SubtitleLocKey:  "SubtitleLocKey",
				SubtitleLocArgs: []string{"Arg1", "Arg2"},
				Body:            "Body",
				LocKey:          "LocKey",
				LocArgs:         []string{"Arg1", "Arg2"},
				Action:          "Action",
				ActionLocKey:    "ActionLocKey",
				SummaryArg:      "SummaryArg",
				SummaryArgCount: 100,
			},
			Badge: 100,
			Sound: Sound{
				Critical: 1,
				Name:     "Sound",
				Volume:   0.1,
			},
			Category:          "Category",
			ThreadID:          "ThreadID",
			MutableContent:    1,
			ContentAvailable:  1,
			TargetContentID:   "TargetContentID",
			InterruptionLevel: InterruptionLevelActive,
			RelevanceScore:    1.0,
			URLArgs:           []string{"Arg1", "Arg2"},
		}, nil),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(n)
	}
}
