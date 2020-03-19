package main

import (
	"log"

	"github.com/dush-t/epirisk/db/models"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// Conn stores information about the current database connection
type Conn struct {
	uri    string
	driver *neo4j.Driver
}

// AddUser adds a new user to the database
func (db Conn) AddUser(phoneNo string, name string) (interface{}, error) {
	driver := *(db.driver)
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Fatal("Failed to connect to database")
		return models.User{}, err
	}
	defer session.Close()

	user, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (user:User) SET user.PhoneNo = $phoneNo SET user.Name = $name SET user.Risk=0.0 RETURN user",
			map[string]interface{}{
				"phoneNo": phoneNo,
				"name":    name,
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

	userDataMap := (user.(neo4j.Node)).Props()
	userEntity := models.BuildUserFromMap(userDataMap)

	return userEntity, nil
}
