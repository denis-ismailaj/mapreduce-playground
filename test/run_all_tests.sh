#!/usr/bin/env bash

# This script by default carries on when a test fails and
# shows a summary in the end.
# To end execution with exit code 1 when a test fails, use --ci-mode.

if [ $# -eq 0 ]; then
  CI_MODE=0
elif [ $# -eq 1 ]; then
  if [ "$1" == '--ci-mode' ]; then
    CI_MODE=1
  else
    echo "Invalid argument: $1. Did you mean --ci-mode?"
    exit 1
  fi
else
  echo "Too many arguments. Expected 1 or 0, found $#."
  exit 1
fi

FAILED_ANY=0

# Helper function to run tests
# It expects the test directory as an argument
function run_test() {
  TEST_IMAGE=$(docker build -q -f "$1"/Dockerfile ..)

  docker run --rm -it "$TEST_IMAGE"
  TEST_RESULT=$?

  # TODO Add option to not retain test images
  # docker rmi "$TEST_IMAGE"

  if [ $TEST_RESULT -ne 0 ]; then
    handle_test_failure
  fi
}

function handle_test_failure() {
  if [ $CI_MODE -eq 0 ]; then
    FAILED_ANY=1
  else
    exit 1
  fi
}

echo '***' Building base test image.
docker build -t mr-playground-test -f ./Dockerfile .. || exit 1

echo '***' Starting test sequence.
run_test wc
run_test indexer
run_test mtiming
run_test rtiming
run_test jobcount
run_test early_exit
run_test crash

printf '\n'

if [ $FAILED_ANY -eq 0 ]; then
    echo '***' PASSED ALL TESTS
    exit 0
else
    echo '***' FAILED SOME TESTS
    exit 1
fi
