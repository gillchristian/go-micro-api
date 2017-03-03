package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Todo struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Content string        `bson:",omitempty" json:"content,omitempty"`
	Title   string        `bson:",omitempty" json:"title,omitempty"`
}

// ErrorWithJSON sets the status code of the response
// and writes a JSON error message to it
func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{\"message\": %q}", message)
}

// ResponseWithJSON sets the status code of the response
// and writes a json to it
func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

func AddTodo(s *mgo.Session, DB DataBaseConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		session := s.Copy()
		defer session.Close()

		var todo Todo
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&todo)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		c := session.DB(DB.Name).C(DB.Collection)

		todo.ID = bson.NewObjectId()

		err = c.Insert(todo)
		if err != nil {
			if mgo.IsDup(err) {
				ErrorWithJSON(w, "Already exists!", http.StatusBadRequest)
				return
			}

			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert todo: ", err)
			return
		}

		respBody, err := json.MarshalIndent(todo, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Location", r.URL.Path+"/"+todo.ID.Hex())
		ResponseWithJSON(w, respBody, http.StatusCreated)
	}
}

func AllTodos(s *mgo.Session, DB DataBaseConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		session := s.Copy()
		defer session.Close()

		c := session.DB(DB.Name).C(DB.Collection)

		var todos []Todo
		err := c.Find(bson.M{}).All(&todos)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all todos: ", err)
			return
		}

		respBody, err := json.MarshalIndent(todos, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func SingleTodo(s *mgo.Session, DB DataBaseConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		session := s.Copy()
		defer session.Close()

		c := session.DB(DB.Name).C(DB.Collection)

		var todo Todo
		id := bson.ObjectIdHex(p.ByName("id"))

		err := c.Find(bson.M{"_id": id}).One(&todo)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get todo: ", err)
			return
		}

		respBody, err := json.MarshalIndent(todo, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func DeleteTodo(s *mgo.Session, DB DataBaseConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		session := s.Copy()
		defer session.Close()

		c := session.DB(DB.Name).C(DB.Collection)

		id := bson.ObjectIdHex(p.ByName("id"))

		err := c.Remove(bson.M{"_id": id})
		if err != nil {
			if err == mgo.ErrNotFound {
				ErrorWithJSON(w, "Book not found", http.StatusNotFound)
				return
			}
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func UpdateTodo(s *mgo.Session, DB DataBaseConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		session := s.Copy()
		defer session.Close()

		var todo Todo
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&todo)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		c := session.DB(DB.Name).C(DB.Collection)

		id := bson.ObjectIdHex(p.ByName("id"))

		err = c.Update(bson.M{"_id": id}, todo)
		if err != nil {
			if err == mgo.ErrNotFound {
				ErrorWithJSON(w, "Book not found", http.StatusNotFound)
				return
			}
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		// send the updated todo Maybe ?
	}
}
