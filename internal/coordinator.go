package internal

import (
	"mapreduce/api"
	"mapreduce/pkg"
	"time"
)

// HandleJobRequest
// an RPC handler for job requests from workers
//
// the RPC argument and reply types are defined in rpc.go.
//
func (c *Coordinator) HandleJobRequest(args *api.JobRequestArgs, reply *api.JobRequestReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	reply.NrReduce = c.nReduce

	for i, job := range c.mapJobs {
		if job.Status == pkg.Unprocessed {
			c.mapJobs[i].Status = pkg.Processing
			c.mapJobs[i].LastStatusUpdate = time.Now()

			reply.Job = c.mapJobs[i]

			return nil
		}

		if job.Status == pkg.Processing && job.IsStale() {
			c.mapJobs[i].Status = pkg.Unprocessed
			c.mapJobs[i].LastStatusUpdate = time.Now()
		}
	}

	// Check if all map jobs are done
	for _, job := range c.mapJobs {
		if job.Status != pkg.Done {
			reply.Job = pkg.Job{Type: pkg.Wait}

			return nil
		}
	}

	// if we're here it means map jobs have finished
	for i, job := range c.reduceJobs {
		if job.Status == pkg.Unprocessed {
			c.reduceJobs[i].Status = pkg.Processing
			c.reduceJobs[i].LastStatusUpdate = time.Now()

			reply.Job = c.reduceJobs[i]

			return nil
		}

		if job.Status == pkg.Processing && job.IsStale() {
			c.reduceJobs[i].Status = pkg.Unprocessed
			c.reduceJobs[i].LastStatusUpdate = time.Now()
		}
	}

	// Check if all reduce jobs are done
	for _, job := range c.reduceJobs {
		if job.Status != pkg.Done {
			reply.Job = pkg.Job{Type: pkg.Wait}

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
func (c *Coordinator) HandleJobFinish(args *api.JobFinishArgs, reply *api.JobFinishReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch args.Job.Type {
	case pkg.Map:
		if c.mapJobs[args.Job.Id].Status == pkg.Done {
			return nil
		}

		for i, output := range args.Outputs {
			c.reduceJobs[i].Inputs = append(c.reduceJobs[i].Inputs, output)
		}

		c.mapJobs[args.Job.Id].Status = pkg.Done
		c.mapJobs[args.Job.Id].LastStatusUpdate = time.Now()
	case pkg.Reduce:
		if c.reduceJobs[args.Job.Id].Status == pkg.Done {
			return nil
		}

		c.reduceJobs[args.Job.Id].Status = pkg.Done
		c.reduceJobs[args.Job.Id].LastStatusUpdate = time.Now()
	}

	return nil
}
