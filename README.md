# MapReduce Playground

## Usage

### Word count example
    // Build plugin
    go build -buildmode=plugin ./test/wc/wc.go 

    // Start coordinator and pass it input files
    go run ./cmd/coordinator test/testdata/pg*.txt

    // Start worker
    go run ./cmd/worker wc.so

## Testing

### Run all tests
    cd test
    ./run_all_tests.sh

### Running a specific test
    // First build the base test image
    docker build -t mr-playground-test -f test/Dockerfile .
    
    // Replace <test> in this command with the test name (e.g. wc)
    docker run --rm -it $(docker build -q -f test/<test>/Dockerfile .)

    
