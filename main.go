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

	// TODO Find a way to make this middleware chaining less ugly
	http.Handle("/sign_in", middleware.LogReq(api.SignInHandler(c)))
	http.Handle("/sign_up", middleware.LogReq(api.SignUpHandler(c)))
	http.Handle("/met_user", middleware.LogReq(middleware.Auth(c, api.MetUserHandler(c))))
	http.Handle("/update_self_healthstatus", middleware.LogReq(middleware.Auth(c, api.UpdateSelfHealthStatus(c))))
	http.Handle("/get_contact_summary", middleware.LogReq(middleware.Auth(c, api.GetContactSummary(c))))

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
