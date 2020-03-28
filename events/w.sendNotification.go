package events

import (
	"fmt"
)

// SendHSChangeNotificationWorker returns a Worker that will send a
// notification to a group of users using FCM
func SendHSChangeNotificationWorker(fcmKey string) Worker {
	return func(e Event) {
		fmt.Println("SendHSChangeNotificationWorker : Sending notification :", e)
	}
}
