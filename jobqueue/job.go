package jobqueue

import (
	"encoding/base64"
	"time"

	"github.com/google/uuid"
)

type JobID string

type Job struct {
	ID          JobID
	Status      Status
	CompletedAt *time.Time
	Input       []byte // CAN filedata
	Result      []byte // Converted DBC filedata
	Err         error
}

type JobSummary struct {
	ID          JobID
	Status      Status
	CompletedAt *time.Time
	Err         error
}

func NewJob(input []byte) *Job {
	return &Job{
		ID:          JobID(generateUUIDBase64()),
		Status:      Pending,
		CompletedAt: nil,
		Input:       input,
		Result:      nil,
		Err:         nil,
	}
}

func JobToSummary(job *Job) JobSummary {
	return JobSummary{
		ID:          job.ID,
		Status:      job.Status,
		CompletedAt: job.CompletedAt,
		Err:         job.Err,
	}
}

func generateUUIDBase64() string {
	uuid := uuid.New()
	return base64.RawURLEncoding.EncodeToString(uuid[:])
}
