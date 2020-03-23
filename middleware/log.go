package middleware

import (
	"log"
	"net/http"
)

// LogReq is a middleware function to log the http requests
// that the application gets
func LogReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		uri := r.URL.String()
		method := r.Method
		log.Println("^^", method, uri)
	})
}
