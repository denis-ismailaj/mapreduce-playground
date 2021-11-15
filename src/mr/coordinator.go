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
		if job.Status == Unprocessed {
			c.mapJobs[i].Status = Processing
			reply.Job = c.mapJobs[i]

			return nil
		}
	}

	// if we're here it means map jobs have finished
	for i, job := range c.reduceJobs {
		if job.Status == Unprocessed {
			c.reduceJobs[i].Status = Processing
			reply.Job = c.reduceJobs[i]

			return nil
		}
	}

	return nil
}

// HandleJobFinish
// an RPC handler for job finish reports from workers
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) HandleJobFinish(args *JobFinishArgs, reply *JobFinishReply) error {
	for i, output := range args.Outputs {
		c.reduceJobs[i] = Job{
			Type:   Reduce,
			Id:     i,
			Status: Unprocessed,
			Inputs: append(c.reduceJobs[i].Inputs, output),
		}
	}

	log.Println(c.reduceJobs)

	return nil
}
