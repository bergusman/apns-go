package apns

import (
	"encoding/json"
	"testing"
)

func TestBuildPayload(t *testing.T) {
	type Payload struct {
		APS    json.RawMessage
		Custom json.RawMessage
	}
	var p Payload

	n := &Notification{
		Payload: BuildPayload(nil, map[string]interface{}{
			"aps": map[string]interface{}{
				"alert": "Hello, Custom!",
			},
			"custom": map[string]interface{}{
				"hello": "Go",
			},
		}),
	}

	bytes, err := json.Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		t.Fatal(err)
	}

	if string(p.APS) != `{"alert":"Hello, Custom!"}` {
		t.Errorf("APS: %v; want: %v", string(p.APS), `{"alert":"Hello, Custom!"}`)
	}
	if string(p.Custom) != `{"hello":"Go"}` {
		t.Errorf("Custom: %v; want: %v", string(p.Custom), `{"hello":"Go"}`)
	}

	n = &Notification{
		Payload: BuildPayload(&APS{
			Alert: "Hello, APS!",
		}, map[string]interface{}{
			"custom": map[string]interface{}{
				"hello": "Go",
			},
		}),
	}

	bytes, err = json.Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		t.Fatal(err)
	}

	if string(p.APS) != `{"alert":"Hello, APS!"}` {
		t.Errorf("APS: %v; want: %v", string(p.APS), `{"alert":"Hello, APS!"}`)
	}
	if string(p.Custom) != `{"hello":"Go"}` {
		t.Errorf("Custom: %v; want: %v", string(p.Custom), `{"hello":"Go"}`)
	}
}
