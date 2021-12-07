package internal

import (
	"fmt"
	"mapreduce/pkg"
	"strconv"
	"time"
)

func (c *Coordinator) jobCreator(inputFiles []string) {
	for {
		c.mu.Lock()

		for c.areCurrentJobsFinished() {
			switch c.currentStage {
			case Start:
				// Create a Map job for each of the input files
				for _, file := range inputFiles {
					id := strconv.Itoa(ihash(file))

					job := pkg.Job{
						Id:               id,
						Status:           pkg.Unprocessed,
						LastStatusUpdate: time.Now(),
						Inputs:           []string{file},
						Type:             pkg.Map,
					}

					c.jobs[id] = &job
				}

				c.currentStage = Map
			case Map:
				// Create nReduce Reduce jobs
				for i := 0; i < c.nReduce; i++ {
					id := fmt.Sprintf("r-%d", i)

					job := pkg.Job{
						Status:           pkg.Unprocessed,
						LastStatusUpdate: time.Now(),
						Id:               id,
						Type:             pkg.Reduce,
						Inputs:           c.mapOutputs[i],
					}

					c.jobs[id] = &job
				}

				c.currentStage = Reduce
			case Reduce:
				// All stages are completed
				c.currentStage = Finished
				c.mu.Unlock()
				return
			}

			// Wait for a job to be marked as done
			c.cond.Wait()
		}

		c.mu.Unlock()
	}
}
