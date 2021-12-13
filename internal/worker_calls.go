package internal

import (
	"log"
	"mapreduce/api"
	"mapreduce/pkg"
)

// JobRequestCall
// makes an RPC call to the coordinator to get a new job
//
// the RPC argument and reply types are defined in rpc.go.
//
func JobRequestCall() (pkg.Job, int) {
	// declare a reply structure.
	reply := api.JobRequestReply{}

	// send the RPC request, wait for the reply.
	call("Coordinator.HandleJobRequest", struct{}{}, &reply)

	return reply.Job, reply.NrReduce
}

// JobFinishCall
// makes an RPC call to the coordinator to report that a job finished
//
// the RPC argument and reply types are defined in rpc.go.
//
func JobFinishCall(job pkg.Job, outputs map[int]string) {
	log.Printf("Finished job %s.", job.Id)

	// declare an argument structure.
	args := api.JobFinishArgs{Outputs: outputs, Job: job}

	// send the RPC request, wait for the reply.
	call("Coordinator.HandleJobFinish", &args, &api.JobFinishReply{})
}
