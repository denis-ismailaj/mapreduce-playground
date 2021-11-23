package mr

// HandleJobRequest
// an RPC handler for job requests from workers
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) HandleJobRequest(args *JobRequestArgs, reply *JobRequestReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	reply.NrReduce = c.nReduce

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
	c.mu.Lock()
	defer c.mu.Unlock()

	switch args.Job.Type {
	case Map:
		c.mapJobs[args.Job.Id].Status = Done

		for i, output := range args.Outputs {
			c.reduceJobs[i] = Job{
				Type:   Reduce,
				Id:     i,
				Status: Unprocessed,
				Inputs: append(c.reduceJobs[i].Inputs, output),
			}
		}
	case Reduce:
		c.reduceJobs[args.Job.Id].Status = Done
	}

	return nil
}
