package config

import (
	"log"
	"os"
)

// FirebaseConf contains data needed to interact with firebase (duh)
type FirebaseConf struct {
	FCMKey string
}

// Firebase is the actual global object that will store the firebase conf
var Firebase FirebaseConf

// Init gets data from the environment and populates the firebase conf
func (f *FirebaseConf) Init() {
	fcmKey, found := os.LookupEnv("FIREBASE_FCM_KEY")
	if !found {
		log.Println(fcmKey)
		log.Fatal("Firebase FCM Key not found. Exiting.")
	}

	(*f).FCMKey = fcmKey
}
