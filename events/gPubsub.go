package events

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
)

// GooglePubSubClient stores data required to push messages to Google PubSub
type GooglePubSubClient struct {
	client *pubsub.Client
	topics map[string]*pubsub.Topic
}

// Init initializes the pubsub client with the underlying data it
// needs in order to publish messages
func (gpsc *GooglePubSubClient) Init(topics []string) {
	ctx := context.Background()
	projectID := "epirisk"

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Println("Failed to create pubsub client:", err)
	}

	(*gpsc).client = client

	for _, topicName := range topics {
		topic, err := client.CreateTopic(ctx, topicName)
		if err != nil {
			log.Println("Failed to create pubsub topic:", err)
		}
		(*gpsc).topics[topicName] = topic
	}

	log.Println("Google PubSub client initialized.")

}

// Publish will publish a message to a given topic
func (gpsc *GooglePubSubClient) Publish(topicName string, data map[string]interface{}) {
	topic, _ := (*gpsc).topics[topicName]

	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("Invalid message format:", err)
		return
	}

	res := topic.Publish(context.Background(), &pubsub.Message{Data: payload})
	log.Println("Message published:", res)
}
