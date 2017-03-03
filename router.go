package main

import (
	"fmt"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

type Handler func(*mgo.Session, DataBaseConfig) httprouter.Handle

type Route struct {
	Method, Pattern string
	Handler         Handler
}

type Routes []Route

func Router(routes *Routes, s *mgo.Session, DB DataBaseConfig) *httprouter.Router {
	fmt.Println("Setting up routes...\n")

	router := httprouter.New()

	for _, route := range *routes {
		fmt.Printf("  - %6v: %v\n", route.Method, route.Pattern)
		router.Handle(route.Method, route.Pattern, route.Handler(s, DB))
	}

	fmt.Println("\nAll routes set up!\n")

	return router
}
