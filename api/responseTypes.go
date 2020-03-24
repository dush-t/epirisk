package api

import "github.com/dush-t/epirisk/db/models"

// UserResponse is used to build the JSON response in http handlers
type UserResponse struct {
	PhoneNo      string
	Name         string
	HealthStatus float64
}

func buildResponseFromUser(u models.User) UserResponse {
	return UserResponse{
		PhoneNo:      u.PhoneNo,
		Name:         u.Name,
		HealthStatus: u.HealthStatus,
	}
}
