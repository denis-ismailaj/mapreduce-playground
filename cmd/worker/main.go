package main

//
// start a worker process, which is implemented
// in mapreduce/mr/main.go. Typically, there will be
// multiple worker processes, talking to one coordinator.
//
// go run main.go wc.so
//
// Please do not change this file.
//

import "mapreduce/internal"
import "os"
import "fmt"

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: mrworker xxx.so\n")
		os.Exit(1)
	}

	mapf, reducef := internal.LoadPlugin(os.Args[1])

	internal.Worker(mapf, reducef)
}
