package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mapreduce/pkg"
	"os"
	"path/filepath"
)

func RunMap(f func(string, string) []KeyValue, job pkg.Job, nReduce int) map[int]string {
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

	err = file.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}

	kva := f(filename, string(content))

	// Partition and write Map output to nReduce files
	// TODO Don't overwrite finished map outputs
	outputs := writeOutput(kva, nReduce, job.Id)

	return outputs
}

func writeOutput(pairs []KeyValue, nReduce int, jobId string) map[int]string {
	var tempFiles = map[int]*os.File{}

	err := os.MkdirAll("out/tmp", os.ModePerm)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Create temporary files for each bucket
	for i := 0; i < nReduce; i++ {
		tempFile, err := ioutil.TempFile("out/tmp", "intermediate")

		if err != nil {
			log.Fatalf(err.Error())
		}

		tempFiles[i] = tempFile
	}

	// Write output to each bucket
	for _, kv := range pairs {
		reduceTaskNr := getReduceTaskNr(kv.Key, nReduce)

		enc := json.NewEncoder(tempFiles[reduceTaskNr])

		err := enc.Encode(&kv)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	fmt.Printf("Writing final outputs for job %s...\n", jobId)

	// Move temp files into final outputs
	var outputs = map[int]string{}

	for i, tempFile := range tempFiles {
		filename := fmt.Sprintf("mr-%s-%d.txt", jobId, i)
		finalPath := filepath.Join("out", filename)

		err := os.Rename(tempFile.Name(), finalPath)
		if err != nil {
			log.Fatalf(err.Error())
		}

		outputs[i] = finalPath
	}

	return outputs
}
