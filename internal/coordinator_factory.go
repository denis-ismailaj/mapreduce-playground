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
	// Ensure output directory exists
	err := os.MkdirAll("out", os.ModePerm)
	if err != nil {
		log.Fatalf(err.Error())
	}

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
	err := rpc.Register(c)
	if err != nil {
		log.Fatalf(err.Error())
	}

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	go func() {
		err := http.Serve(l, nil)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}()
}
