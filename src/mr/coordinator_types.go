package mr

import "sync"

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
	Id     int
	Status JobStatus
	Inputs []string
	Type   JobType
}
