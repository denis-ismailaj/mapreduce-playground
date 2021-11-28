package mr

import (
	"sync"
	"time"
)

type Coordinator struct {
	nReduce      int
	mapJobs      []Job
	reduceJobs   []Job
	lastMapJobId int
	mu           sync.Mutex
}

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
)

type Job struct {
	Id               int
	Status           JobStatus
	LastStatusUpdate time.Time
	Inputs           []string
	Type             JobType
}

func (job Job) isStale() bool {
	maxAge := time.Now().Add(-10 * time.Second)

	return job.LastStatusUpdate.Before(maxAge)
}
