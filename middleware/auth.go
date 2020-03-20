package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dush-t/epirisk/api"
	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/query"
)

// ContextKey because Go throws a warning if I use a string directly
// to access a value by in request context
type ContextKey string

// Auth checks for the Authorization header, decodes
// the token and adds the user to request context.
func Auth(c db.Conn, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("Authorization"), " ")[1]
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

		user, _ := query.GetUser(c, claims.PhoneNo)
		// key := ContextKey("user")
		key := "user"
		ctx := context.WithValue(r.Context(), key, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
