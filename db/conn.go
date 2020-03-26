package db

import (
	"log"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// Conn holds connection information to the database
type Conn struct {
	URI    string
	Driver *neo4j.Driver
}

// QueryContext is the data type meant to be passed in
// neo4j transactions for dynamic queries
type QueryContext map[string]interface{}

// Connect connects to the database
// and store driver information in the connection struct
func (c *Conn) Connect(username string, password string) {
	driver, err := neo4j.NewDriver(c.URI, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	(*c).Driver = &driver
}
