package events

// EventBus will route events to appropriate handlers
type EventBus struct {
	routes  map[string]EventRoute
	channel *EventChan
}

// Register will register a route in the EventBus
func (eb EventBus) Register(er EventRoute) {
	topic := er.Topic().(string)
	eb.routes[topic] = er
}

// Publish takes an event and routes it to the appropriate
// route
func (eb EventBus) Publish(ed Event) {
	route := eb.routes[ed.topic]
	route.Consume(ed)
}

// Init should be used to initialize the EventBus with EventRoutes
// by repeatedly calling Register()
func (eb EventBus) Init(routes []EventRoute) {
	for _, route := range routes {
		eb.Register(route)
	}
}
