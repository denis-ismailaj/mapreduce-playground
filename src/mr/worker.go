package mr

import (
	"os"
	"reflect"
)

// Worker
// main/mrworker.go calls this function.
//
func Worker(
	mapFun func(string, string) []KeyValue,
	reduceFun func(string, []string) string,
) {
	job, nReduce := JobRequestCall()

	if reflect.ValueOf(job).IsZero() {
		// No jobs left to do
		os.Exit(0)
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

	// It finished so go get another job
	Worker(mapFun, reduceFun)
}
