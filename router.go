package main

import (
	"fmt"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

type Route struct {
	Method, Pattern string
	Handler         func(*mgo.Session) httprouter.Handle
}

type Routes []Route

func Router(routes *Routes, s *mgo.Session) *httprouter.Router {
	fmt.Println("Setting up routes...\n")

	router := httprouter.New()

	for _, route := range *routes {
		fmt.Printf("  - %6v: %v\n", route.Method, route.Pattern)
		router.Handle(route.Method, route.Pattern, route.Handler(s))
	}

	fmt.Println("\nAll routes set up!\n")

	return router
}
