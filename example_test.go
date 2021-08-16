package apns_test

import (
	"fmt"
	"log"
	"time"

	"github.com/bergusman/apns-go"
)

func Example() {
	key, err := apns.AuthKeyFromFile("testdata/AuthKey_5MDQ4KLTY7.p8")
	if err != nil {
		log.Fatal(err)
	}

	token := apns.NewToken(key, "5MDQ4KLTY7", "SUPERTEEM1")
	client := apns.NewClient(token, nil)

	n := &apns.Notification{
		DeviceToken: "7c968c83f6fd6de5843c309150ed1a706bc64fcdc42310f66054c0271e67219e",
		Topic:       "com.example.app",
		ID:          "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D",
		Host:        apns.HostDevelopment,
		PushType:    apns.PushTypeAlert,
		Priority:    apns.PriorityHigh,

		Payload: apns.BuildPayload(&apns.APS{
			Alert: apns.Alert{
				Title: "Hello",
			},
			Sound: apns.SoundDefault,
		}, map[string]interface{}{
			"push-id": "EC1BF194-B3B2-424A-89A9-5A918A6E6B5D",
			"type":    "hello",
		}),
	}

	res, err := client.Push(n)
	if err != nil {
		log.Fatal(err)
	}
	if res.Status == apns.Status200 {
		fmt.Println("Successfully sent!")
	} else {
		fmt.Println("Sent failed by reason:", res.Reason)
	}
}

func ExampleNotification() {

}

func ExampleBuildPayload() {

}

func ExampleGenerateBearer() {
	key, err := apns.AuthKeyFromFile("testdata/AuthKey_5MDQ4KLTY7.p8")
	if err != nil {
		log.Fatal(err)
	}

	token, err := apns.GenerateBearer(key, "5MDQ4KLTY7", "SUPERTEEM1", time.Now().Unix())
	if err != nil {
		log.Fatal(err)
	}

	n := &apns.Notification{
		// Notification parameters and payload
	}

	req, err := n.BuildRequest()
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("authorization", "bearer "+token)
	// res, err := http.DefaultClient.Do(req)
}
