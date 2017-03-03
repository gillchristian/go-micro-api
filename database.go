package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
)

// Connect connects to a MongoDB instance
func Connect(url string) (*mgo.Session, error) {
	fmt.Println("Connecting to the DataBase...")

	session, err := mgo.Dial(url)

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected!!!\n")

	return session, nil
}

// EnsureIndex calls EnsureIndex on the microservice's DB Collection
func EnsureIndex(s *mgo.Session, DB DataBaseConfig) {
	session := s.Copy()
	defer session.Close()

	c := session.DB(DB.Name).C(DB.Collection)

	index := mgo.Index{
		Key:        []string{"_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}
