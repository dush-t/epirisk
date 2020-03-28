package events

import (
	"github.com/dush-t/epirisk/db/models"
	"github.com/dush-t/epirisk/util"
)

func generateUsersMetNotif(user1, user2 models.User) (map[string]interface{}, bool) {
	n := map[string]interface{}{
		"title":     "Alert",
		"notifType": "MetRiskyUser",
		"channel":   "MetUser",
	}

	if user1.HealthStatus == user2.HealthStatus {
		return nil, false
	}

	if user1.HealthStatus > user2.HealthStatus {
		n["to"] = []string{user2.RegToken}
		switch user1.HealthStatus {
		case 1.0:
			n["body"] = "You just came close to a user who has been tested positive. Please stay safe."
		case 0.9:
			n["body"] = "You just came close to a user who has been feeling symptoms. Please stay safe."
		}
		return n, true
	} else {
		n, sendNotif := generateUsersMetNotif(user2, user1)
		return n, sendNotif
	}
}

func AnalyseUsersMetEvent(fcmKey string) Worker {
	return func(e Event) {
		event := e.(UsersMetEvent)
		u1 := event.User1
		u2 := event.User2

		notif, sendNotif := generateUsersMetNotif(u1, u2)
		if !sendNotif {
			return
		}

		regToken := notif["to"].([]string)[0]
		util.SendNotification(fcmKey, regToken, notif)
	}
}
