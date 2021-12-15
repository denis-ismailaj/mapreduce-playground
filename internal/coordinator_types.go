package internal

import (
	"mapreduce/pkg"
	"sync"
)

type Coordinator struct {
	nReduce       int
	workerTimeout int
	jobs          map[string]*pkg.Job
	mapOutputs    map[int][]string
	currentStage  Stage
	mu            *sync.Mutex
	cond          *sync.Cond
}

type Stage int64

const (
	Start Stage = iota
	Map
	Reduce
	Finished
)
