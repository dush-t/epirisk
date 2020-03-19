package main

import (
	"log"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func main() {
	// a, _ := testdb("bolt://localhost:7687", "neo4j", "lolmao12345")
	uri := "bolt://localhost:7687"
	driver, _ := neo4j.NewDriver(uri, neo4j.BasicAuth("neo4j", "lolmao12345", ""))
	db := Conn{uri: uri, driver: &driver}
	a, _ := db.AddUser("6377653833", "Dushyant")

	log.Println(reflect.TypeOf(a))
	log.Println("Hello")
	log.Println(a)
}
