package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/movies", app.allMovies)
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthCheck)
	router.HandlerFunc(http.MethodGet, "/movies/:id", app.showMovie)
	router.HandlerFunc(http.MethodPost, "/movies", app.createMovie)
	router.HandlerFunc(http.MethodDelete, "/movies/:id", app.deleteMovie)
	router.HandlerFunc(http.MethodPatch, "/movies/:id", app.updateMovie)

	return router
}
