package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mapreduce/pkg"
	"os"
	"path/filepath"
	"sort"
)

func RunMap(
	f func(string, string) []KeyValue,
	job pkg.Job,
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
	job pkg.Job,
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

	tempFile, _ := ioutil.TempFile("", "temp")

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
		fmt.Fprintf(tempFile, "%v %v\n", kva[i].Key, output)

		i = j
	}

	tempFile.Close()

	outputName := fmt.Sprintf("mr-out-%s.txt", job.Id)
	os.Rename(tempFile.Name(), filepath.Join("out", outputName))
}
