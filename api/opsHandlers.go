package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/models"
	"github.com/dush-t/epirisk/db/query"
	"github.com/dush-t/epirisk/util"
)

// ContextKey because Go throws a warning if I use a string directly
// to access a value by in request context. Sorry being lazy.
type ContextKey string

// MetUserHandler is the handler called at /met_user
func MetUserHandler(c db.Conn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(models.User)

		var reqBody struct {
			PhoneNo   string `json:"phoneNo"`
			TimeSpent int64  `json:"timeSpent"`
		}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		edge, err := query.MetUser(c, user.PhoneNo, reqBody.PhoneNo, reqBody.TimeSpent)
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

// UpdateSelfRisk is the handler called at /update_self_risk
func UpdateSelfRisk(c db.Conn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(models.User)

		var reqBody struct {
			SelfDiagnosedRisk float64 `json:"selfDiagnosedRisk"`
		}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		riskValueAllowed := util.RiskValueIsAllowed(reqBody.SelfDiagnosedRisk)

		if err != nil || !riskValueAllowed {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err = query.MarkSelfInfected(c, user, reqBody.SelfDiagnosedRisk)
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
