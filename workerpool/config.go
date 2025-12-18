package workerpool

import "time"

type Config struct {
	WorkerCount     int
	CleanUpCount    int
	CleanUpInterval time.Duration
}
