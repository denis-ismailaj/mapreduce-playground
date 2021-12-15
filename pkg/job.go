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
	Id               string
	Status           JobStatus
	LastStatusUpdate time.Time
	Inputs           []string
	Type             JobType
}

func (job Job) IsRunningForMoreThan(seconds int) bool {
	if job.Status != Processing {
		return false
	}

	maxAge := time.Now().Add(time.Duration(-seconds) * time.Second)

	return job.LastStatusUpdate.Before(maxAge)
}
