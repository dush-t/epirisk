package init

import (
	"log"
	"os"

	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/events"
)

// Config holds all the information that needs to be passed
// down to functions in order to make the application work.
// Dependency injection ftw.
type Config struct {
	DBConn   db.Conn
	Bus      events.Bus
	Firebase firebaseConfig
}

type firebaseConfig struct {
	FCMKey string
}

// InitializeApp does everything required to start the app
// including connecting to the database and creating an eventbus.
// It returns a config object from environment variables which
// is to be used throughout the app using dependency injection
func InitializeApp() Config {
	uri, username, password := getDBConnectionParams()
	c := db.Conn{URI: uri}
	c.Initialize(username, password)

	bus := events.MakeBus()

	fConfig := getFirebaseConfig()

	return Config{
		DBConn:   c,
		Bus:      bus,
		Firebase: fConfig,
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

func getFirebaseConfig() firebaseConfig {
	fcmKey, found := os.LookupEnv("FIREBASE_FCM_KEY")
	if !found {
		log.Fatal("Firebase FCM key not found in environment. Quitting.")
	}

	return firebaseConfig{
		FCMKey: fcmKey,
	}
}
