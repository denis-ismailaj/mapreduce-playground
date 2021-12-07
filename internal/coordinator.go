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

	// Try to find an unprocessed or stale job to assign to the worker
	for id, job := range c.jobs {
		if job.Status == pkg.Unprocessed || job.IsStale() {
			c.jobs[id].Status = pkg.Processing
			c.jobs[id].LastStatusUpdate = time.Now()

			reply.NrReduce = c.nReduce
			reply.Job = *c.jobs[id]

			return nil
		}
	}

	// Ask worker to check again for open jobs after a while
	reply.Job = pkg.Job{Type: pkg.Wait}

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

	// Dismiss results if the job has already been completed by another worker
	if c.jobs[args.Job.Id].Status == pkg.Done {
		return nil
	}

	// Save intermediate Map output locations
	if args.Job.Type == pkg.Map {
		for i, output := range args.Outputs {
			c.mapOutputs[i] = append(c.mapOutputs[i], output)
		}
	}

	// Mark the job as Done
	c.jobs[args.Job.Id].Status = pkg.Done
	c.jobs[args.Job.Id].LastStatusUpdate = time.Now()

	// Notify job creator
	c.cond.Broadcast()

	return nil
}
