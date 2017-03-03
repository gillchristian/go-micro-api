package main

import (
	"net/http"
)

var routes = Routes{
	Route{
		Method:  http.MethodPost,
		Pattern: "/api/todo",
		Handler: AddTodo,
	},
	Route{
		Method:  http.MethodGet,
		Pattern: "/api/todo",
		Handler: AllTodos,
	},
	Route{
		Method:  http.MethodGet,
		Pattern: "/api/todo/:id",
		Handler: SingleTodo,
	},
	Route{
		Method:  http.MethodPut,
		Pattern: "/api/todo/:id",
		Handler: UpdateTodo,
	},
	Route{
		Method:  http.MethodDelete,
		Pattern: "/api/todo/:id",
		Handler: DeleteTodo,
	},
}
