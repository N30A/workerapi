package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/n30a/workerapi/jobqueue"
)

type JobResponse struct {
	JobID       string  `json:"id"`
	Status      string  `json:"status"`
	CompletedAt *string `json:"completed_at"`
	Error       *string `json:"error"`
}

type JobListResponse struct {
	Count int           `json:"count"`
	Jobs  []JobResponse `json:"jobs"`
}

func (app *App) convertListHandler(writer http.ResponseWriter, request *http.Request) {
	jobs, err := app.Queue.GetJobSummaries()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	response := JobListResponse{
		Count: len(jobs),
		Jobs:  make([]JobResponse, 0, len(jobs)),
	}

	for _, job := range jobs {
		var completedAt *string
		if job.CompletedAt != nil {
			s := job.CompletedAt.Format(time.RFC3339)
			completedAt = &s
		}

		var errMsg *string
		if job.Err != nil {
			s := job.Err.Error()
			errMsg = &s
		}

		response.Jobs = append(response.Jobs, JobResponse{
			JobID:       string(job.ID),
			Status:      job.Status.String(),
			CompletedAt: completedAt,
			Error:       errMsg,
		})
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "failed to encode response", http.StatusInternalServerError)
	}
}

func (app *App) convertHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Content-Type") != "application/octet-stream" {
		http.Error(writer, "unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	data, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "failed to read body", http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	mode := request.URL.Query().Get("mode")
	if mode == "direct" {
		// TODO: do work directly and error handling
		// w.Header().Set("Content-Type", "application/octet-stream")
		time.Sleep(time.Second * 2)
		writer.WriteHeader(http.StatusOK)
		return
	}

	job := jobqueue.NewJob(data)
	app.Queue.Enqueue(job)

	writer.Header().Set("Location", fmt.Sprintf("/convert/%s", job.ID))
	writer.WriteHeader(http.StatusAccepted)
}

func (app *App) convertStatusHandler(writer http.ResponseWriter, request *http.Request) {
	jobID := jobqueue.JobID(chi.URLParam(request, "jobID"))

	job, err := app.Queue.GetJobSummary(jobID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	var completedAt *string
	if job.CompletedAt != nil {
		s := job.CompletedAt.Format(time.RFC3339)
		completedAt = &s
	}

	var errMsg *string
	if job.Err != nil {
		s := job.Err.Error()
		errMsg = &s
	}

	response := JobResponse{
		JobID:       string(job.ID),
		Status:      job.Status.String(),
		CompletedAt: completedAt,
		Error:       errMsg,
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "failed to encode response", http.StatusInternalServerError)
	}

}

func (app *App) convertResultHandler(writer http.ResponseWriter, request *http.Request) {

}

func (app *App) convertCancelHandler(writer http.ResponseWriter, request *http.Request) {

}

func (app *App) convertRetryHandler(writer http.ResponseWriter, request *http.Request) {

}
