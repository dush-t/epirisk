package api

import "github.com/dush-t/epirisk/db/models"

// UserResponse is used to build the JSON response in http handlers
type UserResponse struct {
	PhoneNo           string
	Name              string
	Risk              float64
	SuspectsInfection bool
	Infected          bool
}

func buildResponseFromUser(u models.User) UserResponse {
	return UserResponse{
		PhoneNo:           u.PhoneNo,
		Name:              u.Name,
		Risk:              u.Risk,
		SuspectsInfection: u.SuspectsInfection,
		Infected:          u.Infected,
	}
}
