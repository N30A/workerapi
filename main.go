package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/n30a/workerapi/jobqueue"
	"github.com/n30a/workerapi/workerpool"
)

func main() {
	if err := loadEnvFromFile(".env"); err != nil {
		log.Fatal(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	queue := jobqueue.New(jobqueue.Config{
		BufferSize: 64,
		TimeToLive: time.Minute * 30,
	})

	pool := workerpool.New(queue, workerpool.Config{
		WorkerCount:     5,
		CleanUpCount:    1,
		CleanUpInterval: time.Minute * 5,
	})

	poolCtx, poolCancel := context.WithCancel(context.Background())
	defer poolCancel()

	pool.Start(poolCtx)

	app := NewApp(queue)
	app.RegisterMiddleware()
	app.RegisterRoutes()

	go app.ListenAndServe()

	<-shutdown
	log.Printf("received shutdown signal, waiting %v before shutting down...", shutdownTimeout)

	poolCancel()
	app.Shutdown()
}
