package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dush-t/epirisk/api"
	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/middleware"
)

func main() {
	uri, username, password := getDBConnectionParams()
	c := db.Conn{URI: uri}
	c.Initialize(username, password)

	http.Handle("/sign_in", api.SignInHandler(c))
	http.Handle("/sign_up", api.SignUpHandler(c))
	http.Handle("/met_user", middleware.Auth(c, api.MetUserHandler(c)))
	http.Handle("/update_self_risk", middleware.Auth(c, api.UpdateSelfRisk(c)))

	log.Println("HTTP server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

// getDBConnectionParams gets information needed to connect to the
// database (uri, username, password) from the environment.
func getDBConnectionParams() (string, string, string) {
	connectionURI, uriExists := os.LookupEnv("NEO4J_CONNECTION_URI")
	dbUsername, usernameExists := os.LookupEnv("NEO4J_USERNAME")
	dbPassword, passwordExists := os.LookupEnv("NEO4J_PASSWORD")

	if uriExists && usernameExists && passwordExists {
		return connectionURI, dbUsername, dbPassword
	}

	// Change the default values from here
	return "bolt://localhost:7687", "neo4j", "lolmao12345"
}
