package mr

import "time"

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
			c.mapJobs[i].LastStatusUpdate = time.Now()

			reply.Job = c.mapJobs[i]

			return nil
		}

		if job.Status == Processing && job.isStale() {
			c.mapJobs[i].Status = Unprocessed
			c.mapJobs[i].LastStatusUpdate = time.Now()
		}
	}

	// Check if all map jobs are done
	for _, job := range c.mapJobs {
		if job.Status != Done {
			reply.Job = Job{Type: Wait}

			return nil
		}
	}

	// if we're here it means map jobs have finished
	for i, job := range c.reduceJobs {
		if job.Status == Unprocessed {
			c.reduceJobs[i].Status = Processing
			c.reduceJobs[i].LastStatusUpdate = time.Now()

			reply.Job = c.reduceJobs[i]

			return nil
		}

		if job.Status == Processing && job.isStale() {
			c.reduceJobs[i].Status = Unprocessed
			c.reduceJobs[i].LastStatusUpdate = time.Now()
		}
	}

	// Check if all reduce jobs are done
	for _, job := range c.reduceJobs {
		if job.Status != Done {
			reply.Job = Job{Type: Wait}

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
		if c.mapJobs[args.Job.Id].Status == Done {
			return nil
		}

		for i, output := range args.Outputs {
			c.reduceJobs[i].Inputs = append(c.reduceJobs[i].Inputs, output)
		}

		c.mapJobs[args.Job.Id].Status = Done
		c.mapJobs[args.Job.Id].LastStatusUpdate = time.Now()
	case Reduce:
		if c.reduceJobs[args.Job.Id].Status == Done {
			return nil
		}

		c.reduceJobs[args.Job.Id].Status = Done
		c.reduceJobs[args.Job.Id].LastStatusUpdate = time.Now()
	}

	return nil
}
