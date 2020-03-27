package events

import "log"

// EventBus will route events to appropriate handlers
type EventBus struct {
	routes  map[string]EventRoute
	channel *EventChan
	conf    BusConf
}

// Register will register a route in the EventBus
func (eb *EventBus) Register(er EventRoute) {
	topic := er.Topic().(string)
	(*eb).routes[topic] = er
}

// Publish takes an event and routes it to the appropriate
// route
func (eb *EventBus) Publish(ed Event) {
	route := (*eb).routes[ed.topic]
	route.Consume(ed)
}

// Start should be used to initialize the EventBus with EventRoutes
// by repeatedly calling Register()
func (eb *EventBus) Init(b BusConf, routes []EventRoute) {
	log.Println("Starting EventBus.")
	(*eb).conf = b
	for _, route := range routes {
		(*eb).routes[route.Topic().(string)] = route
	}
}
