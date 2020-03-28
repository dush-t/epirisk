package events

import "github.com/dush-t/epirisk/db/models"

// Event contains information about the Topic of the event (used
// to determine where the event will be sent) and the data of the event
type Event interface{}

// EventChan is a channel for the Event data type
type EventChan chan Event

// Worker is just a function with accepts an Event and does
// something with it
type Worker func(Event)

// WorkerList is a slice of workers
type WorkerList []Worker

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

// Bus is an entity that can route events to different EventRoutes
// based on their topic
type Bus interface {
	// Init is used to build a Bus by registering
	// multiple EventRoutes in it. It's like a constructor
	Init(BusConf, []EventRoute)
	// Register adds a new 'endpoint' to the bus
	Register(EventRoute)
	// Publish pushes a new event to the Bus
	Publish(string, Event)
}

// AppNotification represents the data required to push a notification
// to a device using FCM
type AppNotification struct {
	To        []models.User
	Title     string
	Body      string
	NotifType string
	Channel   string
}

// BusConf contains information needed by the Bus (passed using
// dependency injection)
type BusConf struct {
	Firebase firebaseConf
}

// FirebaseConf contains data needed to interact with firebase (duh)
type firebaseConf struct {
	FCMKey string
}

// UsersMetEvent is an event that stores information about
// a user meeting another user
type UsersMetEvent struct {
	User1     models.User
	User2     models.User
	TimeSpent int64
}
