package mr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)
import "log"
import "net/rpc"
import "hash/fnv"

// KeyValue
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//
// Partitioning function for intermediate Outputs
//
func getReduceTaskNr(key string, nReduce int) int {
	return ihash(key) % nReduce
}

// Worker
// main/mrworker.go calls this function.
//
func Worker(
	mapf func(string, string) []KeyValue,
	reducef func(string, []string) string,
) {
	filename, nReduce := JobRequestCall()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
	}
	file.Close()
	kva := mapf(filename, string(content))

	outputs := writeOutput(kva, nReduce)

	JobFinishCall(filename, outputs)
}

func writeOutput(pairs []KeyValue, nReduce int) map[int]string {
	var outputs = map[int]string{}

	for _, kv := range pairs {

		reduceTaskNr := getReduceTaskNr(kv.Key, nReduce)
		filename := fmt.Sprintf("mr-420-%d.txt", reduceTaskNr)

		fo, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		enc := json.NewEncoder(fo)

		enc.Encode(&kv)

		if err := fo.Close(); err != nil {
			panic(err)
		}

		outputs[reduceTaskNr] = filename
	}

	return outputs
}

// JobRequestCall
// makes an RPC call to the coordinator to get a new job
//
// the RPC argument and reply types are defined in rpc.go.
//
func JobRequestCall() (string, int) {

	// declare an argument structure.
	args := JobRequestArgs{}

	// fill in the argument(s).
	args.X = 99

	// declare a reply structure.
	reply := JobRequestReply{}

	// send the RPC request, wait for the reply.
	call("Coordinator.HandleJobRequest", &args, &reply)

	return reply.Filename, reply.NrReduce
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

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcName string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockName := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockName)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcName, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
