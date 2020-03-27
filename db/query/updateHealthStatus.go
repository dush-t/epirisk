package query

import (
	"log"
	"time"

	"github.com/dush-t/epirisk/constants"
	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/models"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func getQueryForStatus(initHStatus, finalHStatus float64) string {
	query := `MATCH (u0:User {PhoneNo: $u0PhoneNo})
			  WITH u0
			  MATCH (u01:User {PhoneNo: $u0PhoneNo})-[r1:MET]-(u1:User) 
			  WHERE r1.LastMet > $minTime `

	if finalHStatus == 1.0 {
		query += `SET u0:POSITIVE SET u0.TestedPositiveTimestamp = $healthStatusChangeTimestamp`
	}
	if finalHStatus == constants.FeelingSymptomsHealthStatus {
		query += `SET u0:HAS_SYMPTOMS SET u0.HasSymptomsTimestamp = $healthStatusChangeTimestamp`
	}

	if finalHStatus == 0.0 && initHStatus == 1.0 {
		query += `SET u0:CURED SET u0.CuredTimestamp = $healthStatusChangeTimestamp`
	}

	if finalHStatus == 0.0 && initHStatus == 0.9 {
		query += `SET u0:FALSE_ALARM`
	}

	query += ` SET u0.HealthStatus = $healthStatus RETURN u0, u1`
	return query
}

// UpdateHealthStatus marks the infected value to true in the database and
// sets risk for other nodes around the infected node.
func UpdateHealthStatus(d db.Conn, u models.User, initHStatus float64, finalHStatus float64) (map[string]interface{}, error) {
	driver := *(d.Driver)
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return nil, err
	}
	defer session.Close()

	query := getQueryForStatus(initHStatus, finalHStatus)

	data, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			query,
			db.QueryContext{
				"u0PhoneNo":                   u.PhoneNo,
				"healthStatus":                finalHStatus,
				"minTime":                     time.Now().AddDate(0, 0, -21).Unix(),
				"healthStatusChangeTimestamp": time.Now().Unix(),
			},
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		if result.Next() {
			res := make(map[string]interface{})

			contact := models.GetUserFromNode(result.Record().GetByIndex(1).(neo4j.Node))
			firstContactList := []models.User{contact}

			res["user"] = models.GetUserFromNode(result.Record().GetByIndex(0).(neo4j.Node))

			for {
				if result.Next() {
					contact = models.GetUserFromNode(result.Record().GetByIndex(1).(neo4j.Node))
					firstContactList = append(firstContactList, contact)
				} else {
					break
				}
			}
			res["firstContactList"] = firstContactList
			return res, nil
		}

		return nil, result.Err()
	})

	if err != nil {
		log.Fatal("Database Access Error:", err)
		return nil, err
	}

	return data.(map[string]interface{}), nil
}
