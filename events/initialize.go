package events

import (
	"log"
	"os"
)

// MakeBus will return a fully loaded (with EventRoutes)
// Bus. This function won't be dynamic, need to hardcode all
// the EventRoutes here.
func MakeBus(b BusConf) Bus {
	var eb EventBus
	routes := []EventRoute{
		HSPositiveRoute(b),
		HSCuredRoute(b),
		HSFeelingSymptomsRoute(b),
		HSDiedRoute(b),
	}
	eb.Init(b, routes)
	return &eb
}

// GetBusConf will return a BusConf built from the environment
func GetBusConf() BusConf {
	fcmKey, found := os.LookupEnv("FIREBASE_FCM_KEY")
	if !found {
		log.Fatal("Firebase FCM key not found. Exiting.")
	}

	return BusConf{
		Firebase: firebaseConf{
			FCMKey: fcmKey,
		},
	}

}
