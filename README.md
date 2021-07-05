# APNs Provider

HTTP/2 Apple Push Notification service (APNs) provider for Go with token-based connection

Example:

```Go
key, err := apns.AuthKeyFromFile("AuthKey_XXXXXXXXXX.p8")
if err != nil {
	log.Fatal(err)
}

token := apns.NewToken(key, "XXXXXXXXXX", "XXXXXXXXXX")
client := apns.NewClient(token, nil)

n := &apns.Notification{
	DeviceToken: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	Host:        apns.HostDevelopment,
	Topic:       "com.example.app",
	Payload: apns.BuildPayload(&apns.APS{
		Alert: "Hi",
	}, nil),
}

fmt.Println(client.Push(n))
```
