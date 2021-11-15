package mr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

// Worker
// main/mrworker.go calls this function.
//
func Worker(
	mapFun func(string, string) []KeyValue,
	reduceFun func(string, []string) string,
) {
	job, nReduce := JobRequestCall()

	log.Println(job)

	switch job.Type {
	case Map:
		kva := RunMap(mapFun, job)

		outputs := writeOutput(kva, nReduce, job.Id)

		JobFinishCall(job.Inputs[0], outputs)
	case Reduce:
		RunReduce(reduceFun, job)
	}
}

func RunMap(
	f func(string, string) []KeyValue,
	job Job,
) []KeyValue {
	// for map there's only one input
	filename := job.Inputs[0]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	file.Close()
	kva := f(filename, string(content))

	return kva
}

func RunReduce(
	f func(string, []string) string,
	job Job,
) {
	var kva []KeyValue

	for _, filename := range job.Inputs {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("cannot open %v", filename)
		}

		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			kva = append(kva, kv)
		}

		file.Close()
	}

	sort.Sort(ByKey(kva))

	outputName := fmt.Sprintf("mr-out-%d.txt", job.Id)
	outputFile, _ := os.Create(outputName)

	//
	// call Reduce on each distinct key in intermediate[],
	// and print the result to mr-out-0.
	//
	i := 0
	for i < len(kva) {
		j := i + 1
		for j < len(kva) && kva[j].Key == kva[i].Key {
			j++
		}
		var values []string
		for k := i; k < j; k++ {
			values = append(values, kva[k].Value)
		}
		output := f(kva[i].Key, values)

		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(outputFile, "%v %v\n", kva[i].Key, output)

		i = j
	}

	outputFile.Close()
}

// JobRequestCall
// makes an RPC call to the coordinator to get a new job
//
// the RPC argument and reply types are defined in rpc.go.
//
func JobRequestCall() (Job, int) {

	// declare an argument structure.
	args := JobRequestArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := JobRequestReply{}

	// send the RPC request, wait for the reply.
	call("Coordinator.HandleJobRequest", &args, &reply)

	return reply.Job, reply.NrReduce
}

// JobFinishCall
// makes an RPC call to the coordinator to report that a job finished
//
// the RPC argument and reply types are defined in rpc.go.
//
func JobFinishCall(filename string, outputs map[int]string) {

	// declare an argument structure.
	args := JobFinishArgs{Outputs: outputs, Filename: filename}

	// send the RPC request, wait for the reply.
	call("Coordinator.HandleJobFinish", &args, JobFinishReply{})
}
