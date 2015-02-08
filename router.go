package main

import "github.com/julienschmidt/httprouter"

func NewRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/todos", TodoIndex)
	router.POST("/todos", TodoCreate)
	router.GET("/todos/:todoId", TodoShow)
	router.DELETE("/todos/:todoId", TodoDelete)

	router.GET("/token", TokenGet)

	return router
}
