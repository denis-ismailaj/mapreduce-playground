package pkg

import "time"

type JobStatus int64

const (
	Unprocessed JobStatus = iota
	Processing
	Done
)

type JobType int64

const (
	Map JobType = iota
	Reduce
	Wait
)

type Job struct {
	Id               int
	Status           JobStatus
	LastStatusUpdate time.Time
	Inputs           []string
	Type             JobType
}

func (job Job) IsStale() bool {
	maxAge := time.Now().Add(-10 * time.Second)

	return job.LastStatusUpdate.Before(maxAge)
}
