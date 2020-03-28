package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dush-t/epirisk/constants"
	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/models"
	"github.com/dush-t/epirisk/db/query"
	"github.com/dush-t/epirisk/events"
	"github.com/dush-t/epirisk/util"
)

// MetUserHandler is the handler called at /met_user
func MetUserHandler(d db.Conn, b events.Bus) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(models.User)

		var reqBody struct {
			PhoneNo     string `json:"phoneNo"`
			TimeSpent   int64  `json:"timeSpent"`
			MeetingTime int64  `json:"meetingTime"`
		}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := query.MetUser(d, user.PhoneNo, reqBody.PhoneNo, reqBody.TimeSpent, reqBody.MeetingTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var timeSpent int64 = (data["edge"].(models.Edge)).TimeSpent

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		payload := struct {
			TimeSpent int64 `json:"timeSpent"`
		}{TimeSpent: timeSpent}
		json.NewEncoder(w).Encode(payload)

		// Emitting relevant events to the EventBus if the users have just started the meeting
		if reqBody.TimeSpent == 0 {
			var usersMetEvent = events.UsersMetEvent{
				User1:     user,
				User2:     data["userMet"].(models.User),
				TimeSpent: timeSpent,
			}

			log.Println("########################################", b)
			b.Publish(constants.UsersMetRouteName, events.Event(usersMetEvent))
		}
	})
}

// UpdateSelfHealthStatus is the handler called at /update_self_risk
func UpdateSelfHealthStatus(d db.Conn, b events.Bus) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(models.User)

		var reqBody struct {
			HealthStatus float64
		}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		riskValueAllowed := util.HealthStatusValueIsAllowed(reqBody.HealthStatus)

		if err != nil || !riskValueAllowed {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := query.UpdateHealthStatus(d, user, user.HealthStatus, reqBody.HealthStatus)
		if err != nil {
			log.Fatal("Error connecting to the database:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user = data["user"].(models.User)
		firstContacts := data["firstContactList"].([]models.User)

		var n events.AppNotification
		n.NotifType = "HS_CHANGE"
		n.Channel = "HS_CHANGE"
		n.Title = "Alert!"
		n.To = firstContacts

		var route string
		switch reqBody.HealthStatus {
		case constants.FeelingSymptomsHealthStatus:
			n.Body = "A user you met in the last 14 days has started feeling symptoms."
			route = constants.HSFeelingSymptomsChannelName
		case 1.0:
			n.Body = "A user you met in the last 14 days has been tested positive."
			route = constants.HSPositiveChannelName
		case 0.0:
			n.Body = "A user you met in the last 14 days has completely recovered!"
			route = constants.HSCuredChannelName
		}

		statusChangeEvent := struct {
			user         models.User
			notification events.AppNotification
		}{user, n}

		b.Publish(route, events.Event(statusChangeEvent))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		payload := buildResponseFromUser(user)
		json.NewEncoder(w).Encode(payload)
	})
}
