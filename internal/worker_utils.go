package internal

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

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

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
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

func writeOutput(pairs []KeyValue, nReduce int, jobId int) map[int]string {
	var outputs = map[int]string{}

	for _, kv := range pairs {

		reduceTaskNr := getReduceTaskNr(kv.Key, nReduce)
		filename := fmt.Sprintf("mr-%d-%d.txt", jobId, reduceTaskNr)

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

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }