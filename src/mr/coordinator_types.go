package mr

type Coordinator struct {
	nReduce      int
	mapJobs      []MapJob
	reduceJobs   []ReduceJob
	lastMapJobId int
}

type JobStatus int64

const (
	Unprocessed JobStatus = iota
	Processing
	Done
)

type MapJob struct {
	inputPath string
	status    JobStatus
}

type ReduceJob struct {
	taskNumber int
	status     JobStatus
	inputs     []string
}
