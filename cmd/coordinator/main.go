package main

//
// start the coordinator process, which is implemented
// in mapreduce/mr/main.go
//
// go run main.go pg*.txt
//
// Please do not change this file.
//

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"mapreduce/internal"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: mrcoordinator inputfiles...\n")
		os.Exit(1)
	}

	// Try to load from .env file. Otherwise, assume variables are loaded.
	err := godotenv.Load()
	if err == nil {
		log.Println("Loaded environment variables from .env.")
	}

	nReduce, err := strconv.Atoi(os.Getenv("N_REDUCE"))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "N_REDUCE environment variable not found or invalid.\n")
		os.Exit(1)
	}

	m := internal.MakeCoordinator(os.Args[1:], nReduce)
	for m.Done() == false {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
}
