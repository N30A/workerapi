package main

import (
	"github.com/go-chi/chi/v5"
)

func (app *App) RegisterRoutes() {
	app.registerConvertRoutes()
}

func (app *App) registerConvertRoutes() {
	app.Router.Route("/convert", func(router chi.Router) {
		router.Get("/", app.convertListHandler)
		router.Post("/", app.convertHandler)
		router.Get("/{jobID}", app.convertStatusHandler)
		router.Get("/{jobID}/result", app.convertResultHandler)
		router.Delete("/{jobID}", app.convertCancelHandler)
		router.Patch("/{jobID}", app.convertRetryHandler)
	})
}
