package query

import (
	"log"

	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/models"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// MetUser takes a models.User struct and a phoneNo and adds
// an edge between the user corresponding to the struct and
// the one corresponding to the phoneNo
func MetUser(c db.Conn, u1PhoneNo string, u2PhoneNo string, timeSpent int64) (models.Edge, error) {
	driver := *(c.Driver)
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return models.Edge{}, err
	}
	defer session.Close()

	edgeTimeSpent, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			MATCH (u1:User {PhoneNo: $u1PhoneNo})
			MATCH (u2:User {PhoneNo: $u2PhoneNo})
			MERGE (u1)-[r:MET]-(u2)
				ON CREATE SET r.TimeSpent = $timeSpent
				ON MATCH SET r.TimeSpent = r.TimeSpent + $timeSpent
			RETURN r.TimeSpent
			`,
			db.QueryContext{
				"u1PhoneNo": u1PhoneNo,
				"u2PhoneNo": u2PhoneNo,
				"timeSpent": timeSpent,
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
		log.Fatal("Error adding edge in database:", err)
		return models.Edge{}, nil
	}

	edgeEntity := models.GenerateEdgeFromInt(edgeTimeSpent.(int64))

	return edgeEntity, nil
}
