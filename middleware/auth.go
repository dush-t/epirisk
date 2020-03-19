package middleware

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dush-t/epirisk/api"
	"github.com/dush-t/epirisk/db"
)

type ContextKey string

// Auth checks for the Authorization header, decodes
// the token and adds the user to request context.
func Auth(c Conn, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := Split(r.Header.Get("Authorization"), " ")[1]
		claims := &api.Claims{}

		tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("lolmao12345"), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := db.GetUser(c, claims.PhoneNo)
		key ContextKey = "user"
		ctx := context.WithValue(r.Context(), key, user)
		next.serveHTTP(w, r.WithContext(ctx))
	})
}
