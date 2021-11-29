package internal

import (
	"mapreduce/pkg"
	"sync"
)

type Coordinator struct {
	nReduce      int
	mapJobs      []pkg.Job
	reduceJobs   []pkg.Job
	lastMapJobId int
	mu           sync.Mutex
}
