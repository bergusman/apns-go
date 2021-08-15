package apns

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

// See https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/sending_notification_requests_to_apns.

// Use HTTP/2 and TLS 1.2 or later to establish a connection
// between your provider server and one of the following servers:
const (
	// Development server.
	HostDevelopment = "https://api.sandbox.push.apple.com"

	// Development server with 2197 port.
	HostDevelopmentPort2197 = "https://api.sandbox.push.apple.com:2197"

	// Production server.
	HostProduction = "https://api.push.apple.com"

	// Production server with 2197 port.
	HostProductionPort2197 = "https://api.push.apple.com:2197"
)

const (
	// Use the alert push type for notifications that trigger
	// a user interaction—for example, an alert, badge, or sound.
	// If you set this push type, the apns-topic header field
	// must use your app’s bundle ID as the topic.
	//
	// If the notification requires immediate action from the user,
	// set notification priority to 10; otherwise use 5.
	//
	// The alert push type is required on watchOS 6 and later.
	// It is recommended on macOS, iOS, tvOS, and iPadOS.
	PushTypeAlert = "alert"

	// Use the background push type for notifications that deliver content in the background,
	// and don’t trigger any user interactions.
	// If you set this push type, the apns-topic header field must use
	// your app’s bundle ID as the topic. Always use priority 5.
	// Using priority 10 is an error. For more information,
	// see https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/pushing_background_updates_to_your_app.
	//
	// The background push type is required on watchOS 6 and later.
	// It is recommended on macOS, iOS, tvOS, and iPadOS.
	PushTypeBackground = "background"

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

// BuildRequest builds remote notification request for notification.
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

// SetHeaders sets headers of notification for remote notification request.
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
