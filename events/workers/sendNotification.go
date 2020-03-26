package workers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dush-t/epirisk/config"
	"github.com/dush-t/epirisk/db/models"
	"github.com/dush-t/epirisk/events"
)

func sendNotification(fcmKey string, u models.User, n events.AppNotification) {
	data := struct {
		To   string                 `json:"to"`
		Data events.AppNotification `json:"data"`
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

func SendNotificationWorker(c config.Config) events.Worker {
	return func(e events.Event) {

	}
}
