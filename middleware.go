package main

import (
	"github.com/go-chi/chi/v5/middleware"
)

func (app *App) RegisterMiddleware() {
	app.Router.Use(middleware.Heartbeat("/health"))
}
