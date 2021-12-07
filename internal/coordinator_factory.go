package internal

import (
	"log"
	"mapreduce/pkg"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

// MakeCoordinator
// create a Coordinator.
// main/main.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(inputFiles []string, nReduce int) *Coordinator {
	c := Coordinator{
		nReduce:      nReduce,
		mapOutputs:   map[int][]string{},
		currentStage: Start,
		jobs:         map[string]*pkg.Job{},
		mu:           &sync.Mutex{},
	}
	c.cond = sync.NewCond(c.mu)

	// Kick off job creator
	go c.jobCreator(inputFiles)

	c.server()
	return &c
}

//
// start a thread that listens for RPCs from main.go
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
