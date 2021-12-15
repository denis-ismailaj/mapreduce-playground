package internal

import (
	"fmt"
	"log"
	"mapreduce/pkg"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
)

// MakeCoordinator
// create a Coordinator.
// main/main.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(inputFiles []string, nReduce int) *Coordinator {
	workerTimeout, err := strconv.Atoi(os.Getenv("WORKER_TIMEOUT"))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "WORKER_TIMEOUT environment variable not found or invalid.\n")
		os.Exit(1)
	}

	c := Coordinator{
		nReduce:       nReduce,
		workerTimeout: workerTimeout,
		mapOutputs:    map[int][]string{},
		currentStage:  Start,
		jobs:          map[string]*pkg.Job{},
		mu:            &sync.Mutex{},
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

	port := os.Getenv("COORDINATOR_PORT")
	l, e := net.Listen("tcp", ":"+port)
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
