package apns

import (
	"encoding/json"
	"testing"
)

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
			Category:         "Category",
			ThreadID:         "ThreadID",
			MutableContent:   1,
			ContentAvailable: 1,
			TargetContentID:  "TargetContentID",
			URLArgs:          []string{"Arg1", "Arg2"},
		}, nil),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(n)
	}
}
