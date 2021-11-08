package mr

import (
	"log"
)
import "net"
import "os"
import "net/rpc"
import "net/http"

type Coordinator struct {
	Jobs    []Job
	nReduce int
}

type JobStatus int64

const (
	Unprocessed JobStatus = iota
	Processing
	Done
)

type Job struct {
	File   string
	Status JobStatus
}

// HandleJobRequest
// an RPC handler for job requests from workers
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) HandleJobRequest(args *JobRequestArgs, reply *JobRequestReply) error {
	reply.NrReduce = c.nReduce

	// pay attention to potential race condition
	for i, job := range c.Jobs {
		if job.Status == Unprocessed {
			c.Jobs[i].Status = Processing
			reply.Filename = job.File
			return nil
		}
	}

	reply.Filename = ""

	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockName := coordinatorSock()
	os.Remove(sockName)
	l, e := net.Listen("unix", sockName)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// Done
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	for _, job := range c.Jobs {
		if job.Status != Processing {
			return false
		}
	}

	return true
}

// MakeCoordinator
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{nReduce: nReduce}

	for _, file := range files {
		c.Jobs = append(c.Jobs, Job{Status: Unprocessed, File: file})
	}

	c.server()
	return &c
}
