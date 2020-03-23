package query

import (
	"log"
	"math"
	"os"
	"strconv"

	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/models"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// MarkSelfInfected marks the infected value to true in the database and
// sets risk for other nodes around the infected node.
func MarkSelfInfected(c db.Conn, u models.User, healthStatusChange float64) (models.User, error) {
	driver := *(c.Driver)
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return models.User{}, err
	}
	defer session.Close()

	infectionProbabilityValue, _ := os.LookupEnv("INFECTION_PROBABILITY")
	infectionProbability, _ := strconv.ParseFloat(infectionProbabilityValue, 64)

	firstContactRisk := healthStatusChange * math.Pow(infectionProbability, 1)
	secondContactRisk := healthStatusChange * math.Pow(infectionProbability, 2)
	thirdContactRisk := healthStatusChange * math.Pow(infectionProbability, 3)
	fourthContactRisk := healthStatusChange * math.Pow(infectionProbability, 4)

	user, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (u0:User {PhoneNo: $u0PhoneNo})-[r1:MET]-(u1:User) 
			SET u0.Risk = u0.Risk + $healthStatusChange 
			SET u0.HealthStatus = u0.HealthStatus + $healthStatusChange
			SET u1.Risk = u1.Risk + $firstContactRisk * r1.TimeSpent 

			WITH u0, u1

			MATCH (u1:User)-[r2:MET]-(u2:User) 
			WHERE id(u2) <> id(u0) 
			SET u2.Risk = u2.Risk + $secondContactRisk * r2.TimeSpent 

			WITH u0, u1, u2

			MATCH (u2:User)-[r3:MET]-(u3:User) 
			WHERE id(u3) <> id(u1) 
			SET u3.Risk = u3.Risk + $thirdContactRisk * r3.TimeSpent 

			WITH u0, u1, u2, u3

			MATCH (u3:User)-[r4:MET]-(u4:User) 
			WHERE id(u4) <> id(u2) 
			SET u4.Risk = u4.Risk + $fourthContactRisk * r4.TimeSpent 

			RETURN u0
			`,
			db.QueryContext{
				"u0PhoneNo":          u.PhoneNo,
				"healthStatusChange": healthStatusChange,
				"firstContactRisk":   firstContactRisk,
				"secondContactRisk":  secondContactRisk,
				"thirdContactRisk":   thirdContactRisk,
				"fourthContactRisk":  fourthContactRisk,
			},
		)
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().GetByIndex(0), nil
		}

		return nil, result.Err()
	})

	if err != nil {
		log.Fatal("Database Access Error:", err)
		return models.User{}, nil
	}

	userEntity := models.GetUserFromNode(user.(neo4j.Node))

	return userEntity, nil
}
