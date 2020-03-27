package events

import "github.com/dush-t/epirisk/constants"

// HSPositiveRoute will handle all events of a patient being tested positive
func HSPositiveRoute(b BusConf) EventRoute {
	var ner NativeEventRoute
	ner.Init(constants.HSPositiveChannelName)
	ner.RegisterWorker(SendNotificationWorker(b.Firebase.FCMKey))
	return &ner
}

// HSFeelingSymptomsRoute will handle all events of a patient starting to feel symptoms
func HSFeelingSymptomsRoute(b BusConf) EventRoute {
	var ner NativeEventRoute
	ner.Init(constants.HSFeelingSymptomsChannelName)
	ner.RegisterWorker(SendNotificationWorker(b.Firebase.FCMKey))
	return &ner
}

// HSCuredRoute will handle all events of a patient being cured
func HSCuredRoute(b BusConf) EventRoute {
	var ner NativeEventRoute
	ner.Init(constants.HSCuredChannelName)
	return &ner
}

// HSDiedRoute will handle all events of a patient dying
func HSDiedRoute(b BusConf) EventRoute {
	var ner NativeEventRoute
	ner.Init(constants.HSDiedChannelName)
	return &ner
}
