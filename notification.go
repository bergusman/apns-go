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

// Values for Notification.PushType field.
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

	// Use the voip push type for notifications that provide information about
	// an incoming Voice-over-IP (VoIP) call. For more information,
	// see https://developer.apple.com/documentation/pushkit/responding_to_voip_notifications_from_pushkit.
	//
	// If you set this push type, the apns-topic header field must use your app’s bundle ID
	// with .voip appended to the end. If you’re using certificate-based authentication,
	// you must also register the certificate for VoIP services.
	// The topic is then part of the 1.2.840.113635.100.6.3.4 or 1.2.840.113635.100.6.3.6 extension.
	//
	// The voip push type is not available on watchOS.
	// It is recommended on macOS, iOS, tvOS, and iPadOS.
	PushTypeVoIP = "voip"

	// Use the complication push type for notifications
	// that contain update information for a watchOS app’s complications.
	// For more information,
	// see https://developer.apple.com/documentation/clockkit/keeping_your_complications_up_to_date.
	//
	// If you set this push type, the apns-topic header field
	// must use your app’s bundle ID with .complication appended to the end.
	// If you’re using certificate-based authentication,
	// you must also register the certificate for WatchKit services.
	// The topic is then part of the 1.2.840.113635.100.6.3.6 extension.
	//
	// The complication push type is recommended for watchOS and iOS.
	// It is not available on macOS, tvOS, and iPadOS.
	PushTypeComplication = "complication"

	// Use the fileprovider push type to signal changes to a File Provider extension.
	// If you set this push type, the apns-topic header field
	// must use your app’s bundle ID with .pushkit.fileprovider appended to the end.
	// For more information,
	// see https://developer.apple.com/documentation/fileprovider/content_and_change_tracking/tracking_your_file_provider_s_changes/using_push_notifications_to_signal_changes.
	//
	// The fileprovider push type is not available on watchOS.
	// It is recommended on macOS, iOS, tvOS, and iPadOS.
	PushTypeFileprovider = "fileprovider"

	// Use the mdm push type for notifications
	// that tell managed devices to contact the MDM server.
	// If you set this push type, you must use the topic
	// from the UID attribute in the subject
	// of your MDM push certificate.
	// For more information,
	// see https://developer.apple.com/documentation/devicemanagement.
	//
	// The mdm push type is not available on watchOS.
	// It is recommended on macOS, iOS, tvOS, and iPadOS.
	PushTypeMDM = "mdm"
)

// Values for Notification.Priority field.
const (
	// Specify 10 to send the notification immediately.
	PriorityHigh = 10

	// Specify 5 to send the notification based on power considerations on the user’s device.
	PriorityLow = 5
)

type Notification struct {
	// The device token is the hexadecimal bytes that identify the user’s device.
	// Your app receives the bytes for this device token w
	// hen registering for remote notifications;
	// see https://developer.apple.com/documentation/usernotifications/registering_your_app_with_apns.
	DeviceToken string

	Host string

	// A canonical UUID that is the unique ID for the notification.
	// If an error occurs when sending the notification,
	// APNs includes this value when reporting the error to your server.
	// If you omit this header,
	// APNs creates a UUID for you and returns it in its response.
	ID string // header: apns-id

	// The topic for the notification.
	// In general, the topic is your app’s bundle ID/app ID.
	// It can have a suffix based on the type of push notification.
	// If you are using a certificate that supports Pushkit VoIP
	// or watchOS complication notifications,
	// you must include this header with bundle ID of you app
	// and if applicable, the proper suffix.
	// If you are using token-based authentication with APNs,
	// you must include this header with the correct bundle ID
	// and suffix combination.
	Topic string // header: apns-topic

	// The value of this header must accurately reflect the contents
	// of your notification’s payload. If there is a mismatch,
	// or if the header is missing on required systems,
	// APNs may return an error, delay the delivery of the notification,
	// or drop it altogether.
	// Required for watchOS 6 and later; recommended for macOS, iOS, tvOS, and iPadOS.
	PushType string // header: apns-push-type

	// The date at which the notification is no longer valid.
	// This value is a UNIX epoch expressed in seconds (UTC).
	// If the value is nonzero, APNs stores the notification
	// and tries to deliver it at least once,
	// repeating the attempt as needed until the specified date.
	// If the value is 0, APNs attempts to deliver
	// the notification only once and doesn’t store it.
	Expiration string // header: apns-expiration

	// The priority of the notification.
	// If you omit this header, APNs sets the notification priority to 10.
	// Specify 10 to send the notification immediately.
	// Specify 5 to send the notification based on power considerations on the user’s device.
	Priority int // header: apns-priority

	// An identifier you use to coalesce multiple notifications
	// into a single notification for the user.
	// Typically, each notification request causes
	// a new notification to be displayed on the user’s device.
	// When sending the same notification more than once,
	// use the same value in this header to coalesce the requests.
	// The value of this key must not exceed 64 bytes.
	CollapseID string // header: apns-collapse-id

	// The JSON payload with the notification’s content.
	// The JSON payload must not be compressed
	// and is limited to a maximum size of 4 KB (4096 bytes).
	// For a Voice over Internet Protocol (VoIP) notification,
	// the maximum size is 5 KB (5120 bytes).
	Payload interface{}
}

// MarshalJSON marshals the notification payload to JSON.
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

// The path to the device token.
// The value of this header is /3/device/<device_token>.
func (n *Notification) Path() string {
	return "/3/device/" + n.DeviceToken
}

// URL builds full URL of remote notification request.
// If Host field omitted HostProduction will be used.
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

// BuildRequestWithContext builds request with context.
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
