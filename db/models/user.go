package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// User represents the User entity stored in the database
type User struct {
	PhoneNo      string
	Password     string
	Name         string
	Risk         float64
	HealthStatus float64
	Infected     bool
}

type claims struct {
	PhoneNo string `json:"phoneNo"`
	jwt.StandardClaims
}

// GenerateJWT generates a jwt for a user
func (u User) GenerateJWT() (string, error) {
	expirationTime := time.Now().Add(60 * 24 * time.Hour)
	claims := &claims{
		PhoneNo: u.PhoneNo,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	jwtKey := []byte("lolmao12345")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetUserFromNode converts a neo4j node to a User Struct
func GetUserFromNode(n neo4j.Node) User {
	u := n.Props()
	user := User{
		PhoneNo:      u["PhoneNo"].(string),
		Password:     u["Password"].(string),
		Name:         u["Name"].(string),
		Risk:         u["Risk"].(float64),
		HealthStatus: u["HealthStatus"].(float64),
		Infected:     u["Infected"].(bool),
	}

	return user
}
