package mr

// Done
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	for _, job := range c.mapJobs {
		if job.Status != Done {
			return false
		}
	}

	if len(c.reduceJobs) < c.nReduce {
		return false
	}

	for _, job := range c.reduceJobs {
		if job.Status != Done {
			return false
		}
	}

	return true
}
