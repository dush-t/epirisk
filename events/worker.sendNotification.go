package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dush-t/epirisk/db/models"
)

func sendNotification(fcmKey string, u models.User, n AppNotification) {
	data := struct {
		To   string          `json:"to"`
		Data AppNotification `json:"data"`
	}{
		u.RegToken,
		n,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("Could not send notification: invalid data format")
		return
	}

	url := "https://fcm.googleapis.com/fcm/send"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+fcmKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Unable to send notification:", err)
	}
	defer resp.Body.Close()
}

// SendNotificationWorker returns a Worker that will send a
// notification to a group of users using FCM
func SendNotificationWorker(fcmKey string) Worker {
	return func(e Event) {
		fmt.Println("SendNotificationWorker : Sending notification :", e)
	}
}
