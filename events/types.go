package events

// Event contains information about the Topic of the event (used
// to determine where the event will be sent) and the data of the event
type Event struct {
	topic string
	data  interface{}
}

// EventChan is a channel for the Event data type
type EventChan chan Event

// EventRoute is an entity that can consume events by passing them
// to several registered workers
type EventRoute interface {
	// Init is used for any initial setup of the EventRoute
	Init(interface{})

	// Topic returns the topic of the EventRoute
	Topic() interface{}

	// Consume spawns a goroutine for every worker to work
	// on the event
	Consume(Event)

	// RegisterWorker registers a worker function to consume
	// any incoming events in this EventRoute
	RegisterWorker(Worker)
}

// Worker is just a function with accepts an Event and does
// something with it
type Worker func(Event)

// WorkerList is a slice of workers
type WorkerList []Worker

// Bus is an entity that can route events to different EventRoutes
// based on their topic
type Bus interface {
	// Register adds a new 'endpoint' to the bus
	Register(EventRoute)
	// Publish pushes a new event to the Bus
	Publish(Event)
}
