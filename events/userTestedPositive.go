package events

import "github.com/dush-t/epirisk/db/models"

// UserTestedPositive is an event handler for when a user tests positive.
// It sends a notification to all the first contacts of the user and publishes
// a message to the Kafka stream
func UserTestedPositive(u models.User, firstContacts []models.User) {

}
