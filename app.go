package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/n30a/workerapi/jobqueue"
)

const shutdownTimeout = time.Second * 10

type App struct {
	server *http.Server
	Router *chi.Mux
	Queue  *jobqueue.JobQueue
}

func NewApp(queue *jobqueue.JobQueue) *App {
	router := chi.NewRouter()
	server := &http.Server{
		Addr:    net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORT")),
		Handler: router,
	}
	return &App{
		server: server,
		Router: router,
		Queue:  queue,
	}
}

func (app *App) ListenAndServe() {
	log.Printf("listening on %s", net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORT")))
	if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func (app *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	log.Println("server stopped gracefully")
}
