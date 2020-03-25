package events

// NativeEventRoute is used to pass events to functions
// native to this application
type NativeEventRoute struct {
	topic   string
	workers WorkerList
}

// RegisterWorker registers a worker function to consume
// any incoming events in this EventRoute
func (ner *NativeEventRoute) RegisterWorker(w Worker) {
	(*ner).workers = append((*ner).workers, w)
}

// Consume spawns a goroutine for every worker to work
// on the event
func (ner NativeEventRoute) Consume(e Event) {
	for _, w := range ner.workers {
		go w(e)
	}
}

// Topic returns the topic of the EventRoute
func (ner NativeEventRoute) Topic() interface{} {
	return ner.topic
}

// Init is used for any initial setup of the EventRoute
func (ner *NativeEventRoute) Init(t interface{}) {
	(*ner).topic = t.(string)
}
