package query

import (
	"log"

	"github.com/dush-t/epirisk/db"
	"github.com/dush-t/epirisk/db/models"
	"github.com/dush-t/epirisk/util"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// AddUser adds a new user to the database
func AddUser(c db.Conn, phoneNo string, password string, name string) (models.User, error) {
	driver := *(c.Driver)
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Fatal("Failed to connect to database")
		return models.User{}, err
	}
	defer session.Close()

	passwordHash, _ := util.HashPassword(password)

	user, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`
			CREATE (user:User) 
				SET user.PhoneNo = $phoneNo 
				SET user.Password = $password 
				SET user.Name = $name 
				SET user.Risk = 0.0 
				SET user.Infected = false 
				SET user.HealthStatus = 0.0
			RETURN user
			`,
			db.QueryContext{
				"phoneNo":  phoneNo,
				"password": passwordHash,
				"name":     name,
			})
		if err != nil {
			log.Fatal(err)
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
