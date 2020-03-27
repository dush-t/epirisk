package query

import (
	"log"
	"reflect"
	"time"

	"github.com/dush-t/epirisk/constants"
	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/models"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func getQueryForStatus(initHStatus, finalHStatus float64) string {
	query := `MATCH (u0:User {PhoneNo: $u0PhoneNo})-[r1:MET]-(u1:User) 
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
func UpdateHealthStatus(d db.Conn, u models.User, initHStatus float64, finalHStatus float64) (models.User, error) {
	driver := *(d.Driver)
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return models.User{}, err
	}
	defer session.Close()

	query := getQueryForStatus(initHStatus, finalHStatus)

	user, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
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
			log.Println("Keys:", result.Record().Keys())
			nodeList, _ := result.Record().Get("u1")
			log.Println("nodelist type:", reflect.TypeOf(nodeList))
			log.Println(nodeList)
			return result.Record().GetByIndex(0), nil
		}

		return nil, result.Err()
	})

	if err != nil {
		log.Fatal("Database Access Error:", err)
		return models.User{}, nil
	}

	// go events.UserTestedPositive(user, {user})

	userEntity := models.GetUserFromNode(user.(neo4j.Node))

	return userEntity, nil
}
