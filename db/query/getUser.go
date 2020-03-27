package query

import (
	"log"

	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/models"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// GetUser queries the database and gets a user by phoneNo.
func GetUser(d db.Conn, phoneNo string) (models.User, error) {
	driver := *(d.Driver)
	session, err := driver.Session(neo4j.AccessModeRead)
	if err != nil {
		log.Fatal("Failed to connect to database")
		return models.User{}, err
	}
	defer session.Close()

	user, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (user:User {PhoneNo: $phoneNo}) RETURN user",
			db.QueryContext{
				"phoneNo": phoneNo,
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
		log.Fatal(err)
		return models.User{}, err
	}

	userEntity := models.GetUserFromNode(user.(neo4j.Node))

	return userEntity, nil
}
