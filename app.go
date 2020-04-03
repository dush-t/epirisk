package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/events"
	"golang.org/x/tools/godoc/util"
)

// Config holds all the information that needs to be passed
// down to functions in order to make the application work.
// Dependency injection ftw.
type Config struct {
	DBConn db.Conn
	Bus    events.Bus
	BConf  events.BusConf
}

// baseDir returns the path of the root of the project folder
func baseDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}

// InitializeApp does everything required to start the app
// including connecting to the database and creating an eventbus.
// It returns a config object from environment variables which
// is to be used throughout the app using dependency injection
func InitializeApp() Config {
	baseDir := baseDir()
	var conf Config

	// Initialize the environment variables
	util.LoadEnvFromPath(baseDir)

	// Setup gcloud auth
	authFilePath := baseDir + "/gcloudSecrets.json"
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", authFilePath)

	// Connect to database
	uri, username, password := getDBConnectionParams()
	c := db.Conn{URI: uri}
	c.Connect(username, password)
	conf.DBConn = c

	// Initialize Firebase
	var b = events.GetBusConf()

	// Initializer Bus for events
	bus := events.MakeBus(b)
	conf.Bus = bus

	return conf
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
