package jobqueue

import "time"

type Config struct {
	BufferSize int
	TimeToLive time.Duration
}
