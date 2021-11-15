package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

// MakeCoordinator
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{nReduce: nReduce, lastMapJobId: 0}

	for _, file := range files {
		job := Job{
			Id:     c.lastMapJobId,
			Status: Unprocessed,
			Inputs: []string{file},
			Type:   Map,
		}
		c.mapJobs = append(c.mapJobs, job)

		c.lastMapJobId = c.lastMapJobId + 1
	}

	for i := 0; i < nReduce; i++ {
		job := Job{
			Status: Unprocessed,
			Id:     i,
			Type:   Reduce,
		}
		c.reduceJobs = append(c.reduceJobs, job)
	}

	c.server()
	return &c
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
