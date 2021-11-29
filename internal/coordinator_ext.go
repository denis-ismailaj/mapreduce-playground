package internal

import "mapreduce/pkg"

// Done
// main/main.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, job := range c.mapJobs {
		if job.Status != pkg.Done {
			return false
		}
	}

	for _, job := range c.reduceJobs {
		if job.Status != pkg.Done {
			return false
		}
	}

	return true
}
