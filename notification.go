package apns

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

const (
	HostDevelopment         = "https://api.sandbox.push.apple.com"
	HostDevelopmentPort2197 = "https://api.sandbox.push.apple.com:2197"
	HostProduction          = "https://api.push.apple.com"
	HostProductionPort2197  = "https://api.push.apple.com:2197"
)

const (
	PushTypeAlert        = "alert"
	PushTypeBackground   = "background"
	PushTypeVoIP         = "voip"
	PushTypeComplication = "complication"
	PushTypeFileprovider = "fileprovider"
	PushTypeMDM          = "mdm"
)

const (
	PriorityLow  = 5
	PriorityHigh = 10
)

type Notification struct {
	DeviceToken string
	Host        string

	ID         string // header: apns-id
	Topic      string // header: apns-topic
	PushType   string // header: apns-push-type
	Expiration string // header: apns-expiration
	Priority   int    // header: apns-priority
	CollapseID string // header: apns-collapse-id

	Payload interface{}
}

func (n *Notification) MarshalJSON() ([]byte, error) {
	switch v := n.Payload.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return json.Marshal(v)
	}
}

func (n *Notification) Path() string {
	return "/3/device/" + n.DeviceToken
}

func (n *Notification) URL() string {
	host := n.Host
	if host == "" {
		host = HostProduction
	}
	return host + n.Path()
}

func (n *Notification) BuildRequest() (*http.Request, error) {
	return n.BuildRequestWithContext(context.Background())
}

func (n *Notification) BuildRequestWithContext(ctx context.Context) (*http.Request, error) {
	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(n)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", n.URL(), payload)
	if err != nil {
		return nil, err
	}
	n.SetHeaders(req.Header)
	return req, nil
}

func (n *Notification) SetHeaders(h http.Header) {
	if n.ID != "" {
		h.Set("apns-id", n.ID)
	}
	if n.Topic != "" {
		h.Set("apns-topic", n.Topic)
	}
	if n.PushType != "" {
		h.Set("apns-push-type", n.PushType)
	}
	if n.Expiration != "" {
		h.Set("apns-expiration", n.Expiration)
	}
	if n.Priority > 0 {
		h.Set("apns-priority", strconv.Itoa(n.Priority))
	}
	if n.CollapseID != "" {
		h.Set("apns-collapse-id", n.CollapseID)
	}
}
