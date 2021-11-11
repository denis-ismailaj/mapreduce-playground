package mr

import (
	"io/ioutil"
	"log"
	"os"
)

// Worker
// main/mrworker.go calls this function.
//
func Worker(
	mapf func(string, string) []KeyValue,
	reducef func(string, []string) string,
) {
	filename, nReduce, jobId := JobRequestCall()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	file.Close()
	kva := mapf(filename, string(content))

	outputs := writeOutput(kva, nReduce, jobId)

	JobFinishCall(filename, outputs)
}

// JobRequestCall
// makes an RPC call to the coordinator to get a new job
//
// the RPC argument and reply types are defined in rpc.go.
//
func JobRequestCall() (string, int, int) {

	// declare an argument structure.
	args := JobRequestArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := JobRequestReply{}

	// send the RPC request, wait for the reply.
	call("Coordinator.HandleJobRequest", &args, &reply)

	return reply.Filename, reply.NrReduce, reply.JobId
}

// JobFinishCall
// makes an RPC call to the coordinator to report that a job finished
//
// the RPC argument and reply types are defined in rpc.go.
//
func JobFinishCall(filename string, outputs map[int]string) {

	// declare an argument structure.
	args := JobFinishArgs{Outputs: outputs, Filename: filename}

	// send the RPC request, wait for the reply.
	call("Coordinator.HandleJobFinish", &args, JobFinishReply{})
}
