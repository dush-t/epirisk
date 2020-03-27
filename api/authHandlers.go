package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/query"
	"github.com/dush-t/epirisk/events"
	"github.com/dush-t/epirisk/util"
)

// Claims stores the data that will be encoded in a user's JWT
type Claims struct {
	PhoneNo string `json:"phoneNo"`
	jwt.StandardClaims
}

// SignInHandler is the handler function for requests at /sign_in
func SignInHandler(d db.Conn, b events.Bus) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Decode request body
		var data struct {
			PhoneNo  string `json:"phoneNo"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := query.GetUser(d, data.PhoneNo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !util.MatchesWithHash(data.Password, user.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString, err := user.GenerateJWT()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		payload := struct {
			Token string `json:"token"`
		}{Token: tokenString}
		json.NewEncoder(w).Encode(payload)
	})
}

// SignUpHandler is the handler function for requests at /sign_up
func SignUpHandler(d db.Conn, b events.Bus) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			PhoneNo  string `json:"phoneNo"`
			Password string `json:"password"`
			Name     string `json:"name"`
			RegToken string `json:"regToken"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := query.AddUser(d, data.PhoneNo, data.Password, data.Name, data.RegToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error creating user:", err)
			return
		}

		tokenString, err := user.GenerateJWT()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error generating token for user:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		payload := struct {
			Token string `json:"token"`
		}{Token: tokenString}
		json.NewEncoder(w).Encode(payload)
	})
}
