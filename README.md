# MapReduce Playground

## Usage

### Running locally on your machine

    // Build plugin
    go build -buildmode=plugin ./test/wc/wc.go 

    // Start coordinator and pass it input files
    go run ./cmd/coordinator test/testdata/pg*.txt

    // Start worker
    go run ./cmd/worker wc.so

### Running locally with `Docker Compose`

You first need to build the plugin and copy it to the root of the project as `plugin.so`.

Then, you can run:

    docker-compose up --build --scale worker=8

Note that `Go` is picky about plugins, and they need to be built with the same `Go` version as the app.
Here we use the `golang:1.17-bullseye` image, so you should do the same to ensure compatibility.

If you just want to try this project out you can run the following command to have the word count plugin 
be built and copied to where it needs to be.

    docker-compose -f ./hack/build_wc_plugin/compose.yml up --build

## Testing

### Run all tests
    cd test
    ./run_all_tests.sh

### Running a specific test
    // First build the base test image
    docker build -t mr-playground-test -f test/Dockerfile .
    
    // Replace <test> in this command with the test name (e.g. wc)
    docker run --rm -it $(docker build -q -f test/<test>/Dockerfile .)
