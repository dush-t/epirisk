package events

import "github.com/dush-t/epirisk/constants"

// UsersMetRoute will get all the events related to
// meeting users who were feeling symptoms or have been tested
// positive
func UsersMetRoute(b BusConf) EventRoute {
	var ner NativeEventRoute
	ner.Init(constants.UsersMetRouteName)
	ner.RegisterWorker(AnalyseUsersMetEvent(b.Firebase.FCMKey))
	return &ner
}
