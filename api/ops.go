package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/models"
	"github.com/dush-t/epirisk/db/query"
)

// ContextKey because Go throws a warning if I use a string directly
// to access a value by in request context. Sorry being lazy.
type ContextKey string

// MetUserHandler is the handler called at /met_user
func MetUserHandler(c db.Conn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("request context:", r.Context())
		user := r.Context().Value("user").(models.User)
		log.Println("User found in context:", user)
		var data struct {
			PhoneNo   string `json:"phoneNo"`
			TimeSpent int64  `json:"timeSpent"`
		}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		edge, err := query.MetUser(c, user.PhoneNo, data.PhoneNo, data.TimeSpent)
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
