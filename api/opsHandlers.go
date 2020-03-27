package api

import (
	"encoding/json"
	"log"
	"net/http"

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

		edge, err := query.MetUser(d, user.PhoneNo, reqBody.PhoneNo, reqBody.TimeSpent, reqBody.MeetingTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		payload := struct {
			TimeSpent int64 `json:"timeSpent"`
		}{TimeSpent: edge.TimeSpent}
		json.NewEncoder(w).Encode(payload)
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

		user, err = query.UpdateHealthStatus(d, user, user.HealthStatus, reqBody.HealthStatus)
		if err != nil {
			log.Fatal("Error connecting to the database:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		payload := buildResponseFromUser(user)
		json.NewEncoder(w).Encode(payload)
	})
}
