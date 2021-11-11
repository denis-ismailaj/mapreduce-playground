package mr

import (
	"log"
)
import "net"
import "os"
import "net/rpc"
import "net/http"

type Coordinator struct {
	nReduce      int
	mapJobs      []MapJob
	reduceJobs   []ReduceJob
	lastMapJobId int
}

type JobStatus int64

const (
	Unprocessed JobStatus = iota
	Processing
	Done
)

type MapJob struct {
	inputPath string
	status    JobStatus
}

type ReduceJob struct {
	taskNumber int
	status     JobStatus
	inputs     []string
}

// HandleJobRequest
// an RPC handler for job requests from workers
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) HandleJobRequest(args *JobRequestArgs, reply *JobRequestReply) error {
	reply.NrReduce = c.nReduce

	// pay attention to potential race condition
	for i, job := range c.mapJobs {
		if job.status == Unprocessed {
			c.mapJobs[i].status = Processing
			reply.Filename = job.inputPath
			reply.JobId = c.lastMapJobId
			c.lastMapJobId = c.lastMapJobId + 1
			return nil
		}
	}

	reply.Filename = ""

	return nil
}

// HandleJobFinish
// an RPC handler for job finish reports from workers
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) HandleJobFinish(args *JobFinishArgs, reply *JobFinishReply) error {
	for i, output := range args.Outputs {
		c.reduceJobs[i] = ReduceJob{status: Unprocessed, inputs: append(c.reduceJobs[i].inputs, output)}
	}

	log.Println(c.reduceJobs)

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
	for _, job := range c.mapJobs {
		if job.status != Done {
			return false
		}
	}

	if len(c.reduceJobs) < c.nReduce {
		return false
	}

	for _, job := range c.reduceJobs {
		if job.status != Done {
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
	c := Coordinator{nReduce: nReduce, lastMapJobId: 0}

	for _, file := range files {
		c.mapJobs = append(c.mapJobs, MapJob{status: Unprocessed, inputPath: file})
	}

	for i := 0; i < nReduce; i++ {
		c.reduceJobs = append(c.reduceJobs, ReduceJob{status: Unprocessed, taskNumber: i})
	}

	c.server()
	return &c
}
