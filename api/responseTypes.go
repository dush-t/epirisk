package api

import "github.com/dush-t/epirisk/db/models"

// UserResponse is used to build the JSON response in http handlers
type UserResponse struct {
	PhoneNo      string
	Name         string
	Risk         float64
	HealthStatus float64
	Infected     bool
}

// RiskDetailsResponse is used to build the JSON response for score
// details in http handlers
type RiskDetailsResponse struct {
	Risk                float64
	AverageRisk         float64
	AverageRiskGoodness float64
}

func buildResponseFromUser(u models.User) UserResponse {
	return UserResponse{
		PhoneNo:      u.PhoneNo,
		Name:         u.Name,
		Risk:         u.Risk,
		HealthStatus: u.HealthStatus,
		Infected:     u.Infected,
	}
}
