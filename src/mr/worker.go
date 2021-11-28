package mr

import (
	"os"
	"reflect"
	"time"
)

// Worker
// main/mrworker.go calls this function.
//
func Worker(
	mapFun func(string, string) []KeyValue,
	reduceFun func(string, []string) string,
) {
	// Get a new job when finished
	defer Worker(mapFun, reduceFun)

	job, nReduce := JobRequestCall()

	if reflect.ValueOf(job).IsZero() {
		// All jobs are finished
		os.Exit(0)
	}

	if job.Type == Wait {
		// No available jobs at the moment
		time.Sleep(1 * time.Second)
		return
	}

	switch job.Type {
	case Map:
		kva := RunMap(mapFun, job)

		outputs := writeOutput(kva, nReduce, job.Id)

		JobFinishCall(job, outputs)
	case Reduce:
		RunReduce(reduceFun, job)

		JobFinishCall(job, map[int]string{})
	}
}
