package query

import (
	"log"

	"github.com/dush-t/epirisk/constants"
	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/models"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// GetContactSummary will query the database and return information about
// the health statuses of people around the user
func GetContactSummary(d db.Conn, u models.User) (models.ContactSummary, error) {
	driver := *(d.Driver)
	session, err := driver.Session(neo4j.AccessModeRead)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
		return models.ContactSummary{}, err
	}
	defer session.Close()

	summary, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (u0:User {PhoneNo: $phoneNo})
			WITH u0
			MATCH (u0)-[r1:MET]-(u1:User)
			WITH sum(CASE WHEN u1.HealthStatus = 0.9 THEN 1 ELSE 0 END) as firstWSCount,
				 sum(CASE WHEN u1.HealthStatus = 1.0 THEN 1 ELSE 0 END) as firstPCount,
				 sum(CASE WHEN u1.HealthStatus = 0.9 THEN r1.TimeSpent ELSE 0 END) as wsTimeSpent,
				 sum(CASE WHEN u1.HealthStatus = 1.0 THEN r1.TimeSpent ELSE 0 END) as pTimeSpent,
				 u0
			MATCH (u0)-[:MET*2]-(u2:User)
			WHERE NOT((u0)-[:MET]-(u2)) AND u0.PhoneNo <> u2.PhoneNo
			RETURN sum(CASE WHEN u2.HealthStatus = 1.0 THEN 1 ELSE 0 END ) as secondPCount,
				   sum(CASE WHEN u2.HealthStatus = 0.9 THEN 1 ELSE 0 END) as secondWSCount,
				   firstWSCount, firstPCount, wsTimeSpent, pTimeSpent
			`,
			db.QueryContext{
				"phoneNo":      u.PhoneNo,
				"healthStatus": constants.FeelingSymptomsHealthStatus,
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
