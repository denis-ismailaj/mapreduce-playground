package internal

import "mapreduce/pkg"

// Done
// main/main.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.currentStage == Finished
}

func (c *Coordinator) areCurrentJobsFinished() bool {
	if len(c.jobs) == 0 {
		return true
	}

	for _, job := range c.jobs {
		if job.Status != pkg.Done {
			return false
		}
	}

	return true
}
