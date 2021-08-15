package apns

// See https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/generating_a_remote_notification.

// Default system sound.
const SoundDefault = "default"

// Strings that indicates the importance and delivery timing of a notification
const (
	// The system presents the notification immediately,
	// lights up the screen, and can play a sound.
	InterruptionLevelActive = "active"

	// The system presents the notification immediately,
	// lights up the screen, and bypasses the mute switch to play a sound.
	InterruptionLevelCritical = "critical"

	// The system adds the notification to the notification list
	// without lighting up the screen or playing a sound.
	InterruptionLevelPassive = "passive"

	// The system presents the notification immediately, lights up the screen,
	// and can play a sound, but won’t break through system notification controls.
	InterruptionLevelTimeSensitive = "time-senstive"
)

// APS represents Apple-defined remote notification payload keys and their custom values.
type APS struct {
	// String or Alert struct.
	// The information for displaying an alert.
	// A dictionary (Alert) is recommended. If you specify a string,
	// the alert displays your string as the body text.
	Alert interface{} `json:"alert,omitempty"`

	// Int or nil.
	// The number to display in a badge on your app’s icon.
	// Specify 0 to remove the current badge, if any.
	Badge interface{} `json:"badge,omitempty"`

	// String or Sound struct.
	// The name of a sound file in your app’s main bundle
	// or in the Library/Sounds folder of your app’s container directory.
	// Specify the string "default" (SoundDefault) to play the system sound.
	// Use this key for regular notifications.
	// For critical alerts, use the sound dictionary (Sound) instead.
	// For information about how to prepare sounds,
	// see https://developer.apple.com/documentation/usernotifications/unnotificationsound.
	Sound interface{} `json:"sound,omitempty"`

	// An app-specific identifier for grouping related notifications.
	// This value corresponds to the threadIdentifier property in the UNNotificationContent object.
	// See https://developer.apple.com/documentation/usernotifications/unmutablenotificationcontent/1649872-threadidentifier,
	ThreadID string `json:"thread-id,omitempty"`

	// The notification’s type.
	// This string must correspond to the identifier of one of the UNNotificationCategory objects you register at launch time.
	// See https://developer.apple.com/documentation/usernotifications/declaring_your_actionable_notification_types.
	Category string `json:"category,omitempty"`

	// The background notification flag.
	// To perform a silent background update,
	// specify the value 1 and don't include the alert, badge, or sound keys in your payload.
	// See https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/pushing_background_updates_to_your_app.
	ContentAvailable int `json:"content-available,omitempty"`

	// The notification service app extension flag.
	// If the value is 1, the system passes the notification
	// to your notification service app extension before delivery.
	// Use your extension to modify the notification’s content.
	// https://developer.apple.com/documentation/usernotifications/modifying_content_in_newly_delivered_notificationss.
	MutableContent int `json:"mutable-content,omitempty"`

	// The identifier of the window brought forward.
	// The value of this key will be populated on the UNNotificationContent object created from the push payload.
	// Access the value using the UNNotificationContent object's targetContentIdentifier property.
	TargetContentID string `json:"target-content-id,omitempty"`

	// A string that indicates the importance and delivery timing of a notification.
	// The string values “passive”, “active”, “time-senstive”, or “critical”
	// correspond to the UNNotificationInterruptionLevel enumeration cases.
	InterruptionLevel string `json:"interruption-level,omitempty"`

	// Float64 or nil.
	// The relevance score, a number between 0 and 1,
	// that the system uses to sort the notifications from your app.
	// The highest score gets featured in the notification summary.
	// See https://developer.apple.com/documentation/usernotifications/unnotificationcontent/3821031-relevancescore.
	RelevanceScore interface{} `json:"relevance-score,omitempty"`

	URLArgs []string `json:"url-args,omitempty"`
}

// The information for displaying an alert.
type Alert struct {
	// The title of the notification.
	// Apple Watch displays this string in the short look notification interface.
	// Specify a string that is quickly understood by the user.
	Title string `json:"title,omitempty"`

	// Additional information that explains the purpose of the notification.
	Subtitle string `json:"subtitle,omitempty"`

	// The content of the alert message.
	Body string `json:"body,omitempty"`

	// The name of the launch image file to display.
	// If the user chooses to launch your app,
	// the contents of the specified image
	// or storyboard file are displayed instead of your app's normal launch image.
	LaunchImage string `json:"launch-image,omitempty"`

	// The key for a localized title string.
	// Specify this key instead of the title key to retrieve
	// the title from your app’s Localizable.strings files.
	// The value must contain the name of a key in your strings file.
	TitleLocKey string `json:"title-loc-key,omitempty"`

	// An array of strings containing replacement values for variables in your title string.
	// Each %@ character in the string specified by the title-loc-key is replaced by a value from this array.
	// The first item in the array replaces the first instance of the %@ character in the string,
	// the second item replaces the second instance, and so on.
	TitleLocArgs []string `json:"title-loc-args,omitempty"`

	// The key for a localized subtitle string.
	// Use this key, instead of the subtitle key,
	// to retrieve the subtitle from your app's Localizable.strings file.
	// The value must contain the name of a key in your strings file.
	SubtitleLocKey string `json:"subtitle-loc-key,omitempty"`

	// An array of strings containing replacement values for variables in your title string.
	// Each %@ character in the string specified by subtitle-loc-key is replaced by a value from this array.
	// The first item in the array replaces the first instance of the %@ character in the string,
	// the second item replaces the second instance, and so on.
	SubtitleLocArgs []string `json:"subtitle-loc-args,omitempty"`

	// The key for a localized message string.
	// Use this key, instead of the body key,
	// to retrieve the message text from your app's Localizable.strings file.
	// The value must contain the name of a key in your strings file.
	LocKey string `json:"loc-key,omitempty"`

	// An array of strings containing replacement values for variables in your message text.
	// Each %@ character in the string specified by loc-key is replaced by a value from this array.
	// The first item in the array replaces the first instance of the %@ character in the string,
	// the second item replaces the second instance, and so on.
	LocArgs []string `json:"loc-args,omitempty"`

	// The string the notification adds to the category’s summary format string.
	// Deprecated since iOS 15.
	SummaryArg string `json:"summary-arg,omitempty"`

	// The number of items the notification adds to the category’s summary format string.
	// Deprecated since iOS 15.
	SummaryArgCount int `json:"summary-arg-count,omitempty"`

	Action       string `json:"action,omitempty"`
	ActionLocKey string `json:"action-loc-key,omitempty"`
}

// A dictionary that contains sound information for critical alerts.
// For regular notifications, use the sound string instead.
type Sound struct {
	// The critical alert flag. Set to 1 to enable the critical alert.
	Critical int `json:"critical,omitempty"`

	// The name of a sound file in your app’s main bundle
	// or in the Library/Sounds folder of your app’s container directory.
	// Specify the string "default" to play the system sound.
	// For information about how to prepare sounds,
	// see https://developer.apple.com/documentation/usernotifications/unnotificationsound.
	Name string `json:"name,omitempty"`

	// The volume for the critical alert’s sound.
	// Set this to a value between 0 (silent) and 1 (full volume).
	Volume float32 `json:"volume,omitempty"`
}

// BuildPayload builds a remote notification payload from Apple-defined keys aps
// and custom keys custom.
// Returns merged payload from aps keys and custom keys.
// Aps will be choosen as result if it was presented in custom also.
func BuildPayload(aps *APS, custom map[string]interface{}) map[string]interface{} {
	p := make(map[string]interface{})
	for k := range custom {
		p[k] = custom[k]
	}
	if aps != nil {
		p["aps"] = aps
	}
	return p
}
