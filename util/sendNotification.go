package util

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// SendNotification sends a notification to a user with a given firebase regToken
func SendNotification(fcmKey string, regToken string, n map[string]interface{}) {
	log.Println("Inside SendNotification")
	n["to"] = regToken

	payload, err := json.Marshal(n)
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
