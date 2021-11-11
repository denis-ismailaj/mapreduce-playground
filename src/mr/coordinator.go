package mr

import (
	"log"
)

// HandleJobRequest
// an RPC handler for job requests from workers
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) HandleJobRequest(args *JobRequestArgs, reply *JobRequestReply) error {
	reply.NrReduce = c.nReduce

	// pay attention to potential race condition
	for i, job := range c.mapJobs {
		if job.status == Unprocessed {
			c.mapJobs[i].status = Processing
			reply.Filename = job.inputPath
			reply.JobId = c.lastMapJobId
			c.lastMapJobId = c.lastMapJobId + 1
			return nil
		}
	}

	reply.Filename = ""

	return nil
}

// HandleJobFinish
// an RPC handler for job finish reports from workers
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) HandleJobFinish(args *JobFinishArgs, reply *JobFinishReply) error {
	for i, output := range args.Outputs {
		c.reduceJobs[i] = ReduceJob{status: Unprocessed, inputs: append(c.reduceJobs[i].inputs, output)}
	}

	log.Println(c.reduceJobs)

	return nil
}
