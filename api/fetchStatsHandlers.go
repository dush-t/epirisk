package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dush-t/epirisk/config"
	"github.com/dush-t/epirisk/db/models"
	"github.com/dush-t/epirisk/db/query"
)

// GetContactSummary will return a JSON response of the ContactSummary
// of the user
func GetContactSummary(c config.Config) http.Handler {
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
