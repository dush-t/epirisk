package query

import (
	"log"

	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/models"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// GetContactSummary will query the database and return information about
// the health statuses of people around the user
func GetContactSummary(c db.Conn, u models.User) (models.ContactSummary, error) {
	driver := *(c.Driver)
	session, err := driver.Session(neo4j.AccessModeRead)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
		return models.ContactSummary{}, err
	}
	defer session.Close()

	summary, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (u10:User {PhoneNo: $phoneNo})-[r11:MET]-(u11:User)-[r12:MET]-(u12:User)
			WHERE id(u10) <> id(u12) 
			AND u11.HealthStatus = 0.8 AND u12.HealthStatus = 0.8

			WITH u10

			MATCH (u10)-[r21:MET]-(u21:User)-[r22:MET]-(u22:User)
			WHERE id(u10) <> id(u22)
			AND u21.HealthStatus = 1.0 AND u22.HealthStatus = 1.0
			
			RETURN
				count(u11) as firstWithSymptoms
				count(u12) as secondWithSymptoms
				count(u21) as firstPositive
				count(u22) as secondPositive
				sum(r11.TimeSpent) as firstWithSymptomsTimeSpent
				sum(r21.TimeSpent) as firstPositiveTimeSpent
			`,
			db.QueryContext{
				"phoneNo": u.PhoneNo,
			},
		)

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record(), nil
		}

		return nil, result.Err()
	})

	if err != nil {
		log.Fatal("Error fetching data from database:", err)
		return models.ContactSummary{}, nil
	}

	contactSummaryEntity := models.BuildContactSummary(summary.(neo4j.Record))

	return contactSummaryEntity, nil

}
