package models

// User represents the User entity stored in the database
type User struct {
	PhoneNo string
	Name    string
	Risk    float64
}

// BuildUserFromMap accepts a map containing user data
// and returns a User struct. It's just a utility function.
func BuildUserFromMap(u map[string]interface{}) User {
	user := User{
		PhoneNo: u["PhoneNo"].(string),
		Name:    u["Name"].(string),
		Risk:    u["Risk"].(float64),
	}

	return user
}
