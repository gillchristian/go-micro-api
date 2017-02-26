package main

import (
	"log"
	"net/http"
)

var DB = DataBaseConfig{"test_mgo_four", "Todo"}

func main() {
	session, err := Connect("localhost:27017")
	// this defer does not work
	// maybe connect on request
	// see https://goo.gl/25HDes
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	EnsureIndex(session)

	router := Router(&routes, session)

	log.Fatal(http.ListenAndServe(":8080", router))
}
