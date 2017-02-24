package main

import (
	"fmt"
	// "strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Todo struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Content string        `json:"content"`
}

func Connect(url string) (*mgo.Session, error) {
	fmt.Println("Connecting to the DataBase...")

	session, err := mgo.Dial(url)

	if err != nil {
		return nil, err
	}

	fmt.Println("CONNECTED!!!\n")

	return session, nil
}

func EnsureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB(DataBaseName).C(CollectionName)

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