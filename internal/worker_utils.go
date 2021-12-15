package internal

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/rpc"
	"os"
)

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcName string, args interface{}, reply interface{}) bool {
	coordinatorHost := os.Getenv("COORDINATOR_HOST")
	coordinatorPort := os.Getenv("COORDINATOR_PORT")
	coordinatorAddress := fmt.Sprintf("%s:%s", coordinatorHost, coordinatorPort)

	c, err := rpc.DialHTTP("tcp", coordinatorAddress)
	if err != nil {
		return true
	}

	defer func(c *rpc.Client) {
		err := c.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(c)

	err = c.Call(rpcName, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

func ihash(key string) int {
	h := fnv.New32a()
	_, err := h.Write([]byte(key))
	if err != nil {
		log.Fatalf(err.Error())
	}
	return int(h.Sum32() & 0x7fffffff)
}

//
// Partitioning function for intermediate Outputs
//
func getReduceTaskNr(key string, nReduce int) int {
	return ihash(key) % nReduce
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }
