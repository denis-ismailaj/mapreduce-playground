package internal

import (
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
		os.Exit(0)
	}

	switch job.Type {
	case pkg.Wait:
		// No available jobs at the moment
		// Wait a bit before checking again
		time.Sleep(1 * time.Second)
	case pkg.Map:
		// Run Map and get the key value pairs
		kva := RunMap(mapFun, job)

		// Partition and write Map output to nReduce files
		// TODO Don't overwrite finished map outputs
		outputs := writeOutput(kva, nReduce, job.Id)

		// Report back to coordinator with the finished output files
		JobFinishCall(job, outputs)
	case pkg.Reduce:
		// Run Reduce function and output final files
		RunReduce(reduceFun, job)

		// Report job finish to coordinator
		JobFinishCall(job, map[int]string{})
	}
}
