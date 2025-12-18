package jobqueue

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type JobQueue struct {
	store      map[JobID]*Job
	channel    chan *Job
	mutex      sync.Mutex
	timeToLive time.Duration
}

func New(config Config) *JobQueue {
	return &JobQueue{
		store:      make(map[JobID]*Job),
		channel:    make(chan *Job, config.BufferSize),
		timeToLive: config.TimeToLive,
	}
}

func (queue *JobQueue) PendingJobs() <-chan *Job {
	return queue.channel
}

func (queue *JobQueue) Enqueue(job *Job) {
	queue.mutex.Lock()
	queue.store[job.ID] = job
	queue.mutex.Unlock()

	queue.channel <- job
}

func (queue *JobQueue) MarkProcessing(jobID JobID) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	job, ok := queue.store[jobID]
	if !ok {
		return
	}

	job.Status = Processing
}

func (queue *JobQueue) MarkCompleted(jobID JobID, result []byte) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	job, ok := queue.store[jobID]
	if !ok {
		return
	}

	now := time.Now()
	job.Status = Completed
	job.CompletedAt = &now
	job.Result = result
	job.Err = nil
	job.Input = nil
}

func (queue *JobQueue) MarkFailed(jobID JobID, err error) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	job, ok := queue.store[jobID]
	if !ok {
		return
	}

	now := time.Now()
	job.Status = Failed
	job.CompletedAt = &now
	job.Err = err
}

func (queue *JobQueue) CleanUp() {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	for _, job := range queue.store {
		if job.CompletedAt == nil {
			continue
		}

		if job.Status != Completed && job.Status != Failed {
			continue
		}

		now := time.Now()
		if now.Sub(*job.CompletedAt) > queue.timeToLive {
			delete(queue.store, job.ID)
		}
	}
}

func (queue *JobQueue) GetJobSummary(jobID JobID) (JobSummary, error) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	job, ok := queue.store[jobID]
	if !ok {
		return JobSummary{}, fmt.Errorf("job with ID %s was not found", jobID)
	}

	summary := JobSummary{
		ID:          job.ID,
		Status:      job.Status,
		CompletedAt: job.CompletedAt,
		Err:         job.Err,
	}

	return summary, nil
}

func (queue *JobQueue) GetJobSummaries() ([]JobSummary, error) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	jobSummaries := make([]JobSummary, 0, len(queue.store))
	for _, job := range queue.store {
		jobSummaries = append(jobSummaries, JobSummary{
			ID:          job.ID,
			Status:      job.Status,
			CompletedAt: job.CompletedAt,
			Err:         job.Err,
		})
	}

	if len(jobSummaries) == 0 {
		return nil, errors.New("no jobs were found")
	}

	return jobSummaries, nil
}
