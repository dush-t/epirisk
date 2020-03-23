package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/models"
	"github.com/dush-t/epirisk/db/query"
)

// GetCurrentRisk returns the user's current score and it's equivalencies
func GetCurrentRisk(c db.Conn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(models.User)
		risk := user.Risk
		averageRiskValue, _ := os.LookupEnv("AVERAGE_RISK")
		averageRiskGoodnessValue, _ := os.LookupEnv("AVERAGE_RISK_GOODNESS")
		averageRisk, _ := strconv.ParseFloat(averageRiskValue, 64)
		averageRiskGoodness, _ := strconv.ParseFloat(averageRiskGoodnessValue, 64)

		payload := RiskDetailsResponse{
			Risk:                risk,
			AverageRisk:         averageRisk,
			AverageRiskGoodness: averageRiskGoodness,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	})
}

// GetContactSummary will return a JSON response of the ContactSummary
// of the user
func GetContactSummary(c db.Conn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(models.User)

		contactSummary, err := query.GetContactSummary(c, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		payload, err := json.Marshal(contactSummary)
		if err != nil {
			log.Fatal("Unable to convert ContactSummary to JSON:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	})
}
