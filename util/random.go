package util

import (
	"strconv"
	"time"
)

// GenerateAnonID takes a phone number and generates an ID
// thats almost guaranteed to be unique.
func GenerateAnonID(phoneNo string) string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	return timestamp + phoneNo[2:6]
}
