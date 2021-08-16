package apns

import (
	"encoding/json"
	"io"
	"net/http"
)

// See https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/handling_notification_responses_from_apns.

// The possible values in the :status header of the response.
const (
	// Success.
	Status200 = 200

	// Bad request.
	Status400 = 400

	// There was an error with the certificate or with the provider’s authentication token.
	Status403 = 403

	// The request contained an invalid :path value.
	Status404 = 404

	// The request used an invalid :method value. Only POST requests are supported.
	Status405 = 405

	// The device token is no longer active for the topic.
	Status410 = 410

	// The notification payload was too large.
	Status413 = 413

	// The server received too many requests for the same device token.
	Status429 = 429

	// Internal server error.
	Status500 = 500

	// The server is shutting down and unavailable.
	Status503 = 503
)

// The possible error codes included in the reason key of a response’s JSON payload.
const (
	// Status 400. The collapse identifier exceeds the maximum allowed size.
	// The same as ReasonInvalidCollapseId but APNs returns InvalidCollapseId instead of BadCollapseId.
	ReasonBadCollapseId = "BadCollapseId"

	// Status 400. The collapse identifier exceeds the maximum allowed size.
	// The same as ReasonBadCollapseId but APNs returns InvalidCollapseId instead of BadCollapseId.
	ReasonInvalidCollapseId = "InvalidCollapseId"

	// Status 400. The specified device token is invalid.
	// Verify that the request contains a valid token and that the token matches the environment.
	ReasonBadDeviceToken = "BadDeviceToken"

	// Status 400. The apns-expiration value is invalid.
	ReasonBadExpirationDate = "BadExpirationDate"

	// Status 400. The apns-id value is invalid.
	ReasonBadMessageId = "BadMessageId"

	// Status 400. The apns-priority value is invalid.
	ReasonBadPriority = "BadPriority"

	// Status 400. The apns-topic value is invalid.
	ReasonBadTopic = "BadTopic"

	// Status 400. The device token doesn’t match the specified topic.
	ReasonDeviceTokenNotForTopic = "DeviceTokenNotForTopic"

	// Status 400. One or more headers are repeated.
	ReasonDuplicateHeaders = "DuplicateHeaders"

	// Status 400. Idle timeout.
	ReasonIdleTimeout = "IdleTimeout"

	// Status 400. The apns-push-type value is invalid.
	ReasonInvalidPushType = "InvalidPushType"

	// Status 400. The device token isn’t specified in the request :path.
	// Verify that the :path header contains the device token.
	ReasonMissingDeviceToken = "MissingDeviceToken"

	// Status 400. The apns-topic header of the request isn’t specified and is required.
	// The apns-topic header is mandatory when the client is connected
	// using a certificate that supports multiple topics.
	ReasonMissingTopic = "MissingTopic"

	// Status 400. The message payload is empty.
	ReasonPayloadEmpty = "PayloadEmpty"

	// Status 400. Pushing to this topic is not allowed.
	ReasonTopicDisallowed = "TopicDisallowed"

	// Status 403. The certificate is invalid.
	ReasonBadCertificate = "BadCertificate"

	// Status 403. The client certificate is for the wrong environment.
	ReasonBadCertificateEnvironment = "BadCertificateEnvironment"

	// Status 403. The provider token is stale and a new token should be generated.
	ReasonExpiredProviderToken = "ExpiredProviderToken"

	// Status 403. The specified action is not allowed.
	ReasonForbidden = "Forbidden"

	// Status 403. The provider token is not valid, or the token signature can't be verified.
	ReasonInvalidProviderToken = "InvalidProviderToken"

	// Status 403. No provider certificate was used to connect to APNs,
	// and the authorization header is missing or no provider token is specified.
	ReasonMissingProviderToken = "MissingProviderToken"

	// Status 404. The request contained an invalid :path value.
	ReasonBadPath = "BadPath"

	// Status 405. The specified :method value isn’t POST.
	ReasonMethodNotAllowed = "MethodNotAllowed"

	// Status 410. The device token is inactive for the specified topic.
	// There is no need to send further pushes to the same device token,
	// unless your application retrieves the same device token,
	// see https://developer.apple.com/documentation/usernotifications/registering_your_app_with_apns.
	ReasonUnregistered = "Unregistered"

	// Status 413. The message payload is too large.
	// For information about the allowed payload size,
	// see https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/sending_notification_requests_to_apns#2947607.
	ReasonPayloadTooLarge = "PayloadTooLarge"

	// Status 429. The provider’s authentication token is being updated too often.
	// Update the authentication token no more than once every 20 minutes.
	ReasonTooManyProviderTokenUpdates = "TooManyProviderTokenUpdates"

	// Status 429. Too many requests were made consecutively to the same device token.
	ReasonTooManyRequests = "TooManyRequests"

	// Status 500. An internal server error occurred.
	ReasonInternalServerError = "InternalServerError"

	// Status 503. The service is unavailable.
	ReasonServiceUnavailable = "ServiceUnavailable"

	// Status 503. The APNs server is shutting down.
	ReasonShutdown = "Shutdown"
)

// Response for APNs request.
// Each response contains a header with fields indicating the status of the response.
// If the request was successful, the body of the response is empty.
// If an error occurred, the body contains a JSON dictionary
// with additional information about the error.
type Response struct {
	// The same value found in the apns-id field of the request’s header.
	// Use this value to identify the notification.
	// If you don’t specify an apns-id field in your request,
	// APNs creates a new UUID and returns it in this header.
	Id string // header: apns-id

	// The HTTP status code.
	Status int // header: :status

	// The error code (specified as a string) indicating the reason for the failure.
	Reason string `reason:"alert,omitempty"`

	// The time, represented in milliseconds since Epoch,
	// at which APNs confirmed the token was no longer valid for the topic.
	// This key is included only when the error in the :status field is 410.
	Timestamp int64 `timestamp:"alert,omitempty"`
}

// ParseResponse parses HTTP response r from APNs request.
// Parses request header and JSON body.
func ParseResponse(r *http.Response) (*Response, error) {
	defer r.Body.Close()
	res := &Response{
		Id:     r.Header.Get("apns-id"),
		Status: r.StatusCode,
	}
	if err := json.NewDecoder(r.Body).Decode(res); err != nil && err != io.EOF {
		return nil, err
	}
	return res, nil
}
