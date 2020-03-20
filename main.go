package main

import (
	"log"
	"net/http"

	"github.com/dush-t/epirisk/api"
	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/middleware"
)

func main() {
	c := db.Conn{URI: "bolt://localhost:7687"}
	c.Initialize("neo4j", "lolmao12345")

	http.Handle("/sign_in", api.SignInHandler(c))
	http.Handle("/sign_up", api.SignUpHandler(c))
	http.Handle("/met_user", middleware.Auth(c, api.MetUserHandler(c)))

	log.Println("HTTP server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
