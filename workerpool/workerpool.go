package workerpool

import (
	"context"
	"math/rand"
	"time"

	"github.com/n30a/workerapi/jobqueue"
)

type WorkerPool struct {
	jobQueue *jobqueue.JobQueue
	config   Config
}

func New(jobQueue *jobqueue.JobQueue, config Config) *WorkerPool {
	return &WorkerPool{
		jobQueue: jobQueue,
		config:   config,
	}
}

func (workerpool *WorkerPool) Start(ctx context.Context) {
	for range workerpool.config.WorkerCount {
		go worker(workerpool.jobQueue, ctx)
	}

	for range workerpool.config.CleanUpCount {
		go cleanUpWorker(workerpool.jobQueue, workerpool.config.CleanUpInterval, ctx)
	}
}

func worker(queue *jobqueue.JobQueue, ctx context.Context) {
	for {
		select {
		case job := <-queue.PendingJobs():

			queue.MarkProcessing(job.ID)
			// Simulera arbete så länge
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

			queue.MarkCompleted(job.ID, []byte("simulera data"))

		case <-ctx.Done():
			return
		}

	}
}

func cleanUpWorker(queue *jobqueue.JobQueue, cleanUpInterval time.Duration, ctx context.Context) {
	ticker := time.NewTicker(cleanUpInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			queue.CleanUp()
		}
	}
}
