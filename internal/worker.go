package internal

import (
	"fmt"
	"log"
	"mapreduce/pkg"
	"os"
	"reflect"
	"time"
)

// Worker
// main/main.go calls this function.
//
func Worker(
	mapFun func(string, string) []KeyValue,
	reduceFun func(string, []string) string,
) {
	// Rerun worker when finished
	defer Worker(mapFun, reduceFun)

	// Request a job from the coordinator
	job, nReduce := JobRequestCall()

	// Exit if there are no more jobs left
	if reflect.ValueOf(job).IsZero() {
		fmt.Println("No more jobs to do. Exiting...")
		os.Exit(0)
	}

	switch job.Type {
	case pkg.Wait:
		// No available jobs at the moment
		// Wait a bit before checking again
		log.Println("No jobs available. Retrying...")
		time.Sleep(1 * time.Second)
	case pkg.Map:
		log.Printf("Got assigned Map job with id %s.", job.Id)

		// Run Map and get the output files back
		outputs := RunMap(mapFun, job, nReduce)

		// Report back to coordinator with the finished output files
		JobFinishCall(job, outputs)
	case pkg.Reduce:
		log.Printf("Got assigned Reduce job with id %s.", job.Id)

		// Run Reduce function and output final files
		RunReduce(reduceFun, job)

		// Report job finish to coordinator
		JobFinishCall(job, map[int]string{})
	}
}
