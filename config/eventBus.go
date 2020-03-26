package config

import "github.com/dush-t/epirisk/events"

// MakeBus will return a fully loaded (with EventRoutes)
// Bus. This function won't be dynamic, need to hardcode all
// the EventRoutes here.
func MakeBus(c Config) events.Bus {
	var eb events.EventBus
	return eb
}
