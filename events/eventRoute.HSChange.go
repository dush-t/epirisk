package events

// HSPositiveRoute will handle all events of a patient being tested positive
func HSPositiveRoute(b BusConf) EventRoute {
	var ner NativeEventRoute
	ner.Init("HS_POSITIVE")
	ner.RegisterWorker(SendNotificationWorker(b.Firebase.FCMKey))
	return &ner
}

// HSFeelingSymptomsRoute will handle all events of a patient starting to feel symptoms
func HSFeelingSymptomsRoute(b BusConf) EventRoute {
	var ner NativeEventRoute
	ner.Init("HS_FEELING_SYMPTOMS")
	ner.RegisterWorker(SendNotificationWorker(b.Firebase.FCMKey))
	return &ner
}

// HSCuredRoute will handle all events of a patient being cured
func HSCuredRoute(b BusConf) EventRoute {
	var ner NativeEventRoute
	ner.Init("HS_CURED")
	return &ner
}

// HSDiedRoute will handle all events of a patient dying
func HSDiedRoute(b BusConf) EventRoute {
	var ner NativeEventRoute
	ner.Init("HS_FEELING_SYMPTOMS")
	return &ner
}
