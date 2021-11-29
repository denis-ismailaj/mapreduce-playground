package api

//
// RPC definitions.
//

import (
	"mapreduce/pkg"
)

type JobRequestArgs struct {
	X int
}

type JobRequestReply struct {
	Job      pkg.Job
	NrReduce int
}

type JobFinishArgs struct {
	Job     pkg.Job
	Outputs map[int]string
}

type JobFinishReply struct{}
