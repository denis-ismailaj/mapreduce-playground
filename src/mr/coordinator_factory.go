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
		c.mapJobs = append(c.mapJobs, MapJob{status: Unprocessed, inputPath: file})
	}

	for i := 0; i < nReduce; i++ {
		c.reduceJobs = append(c.reduceJobs, ReduceJob{status: Unprocessed, taskNumber: i})
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
