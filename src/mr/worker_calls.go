package mr

// JobRequestCall
// makes an RPC call to the coordinator to get a new job
//
// the RPC argument and reply types are defined in rpc.go.
//
func JobRequestCall() (Job, int) {
	// declare a reply structure.
	reply := JobRequestReply{}

	// send the RPC request, wait for the reply.
	call("Coordinator.HandleJobRequest", struct{}{}, &reply)

	return reply.Job, reply.NrReduce
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
	call("Coordinator.HandleJobFinish", &args, &JobFinishReply{})
}
