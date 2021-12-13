package main

//
// start the coordinator process, which is implemented
// in mapreduce/mr/main.go
//
// go run main.go pg*.txt
//
// Please do not change this file.
//

import "mapreduce/internal"
import "time"
import "os"
import "fmt"

func main() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: mrcoordinator inputfiles...\n")
		os.Exit(1)
	}

	m := internal.MakeCoordinator(os.Args[1:], 10)
	for m.Done() == false {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
}
