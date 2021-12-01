#!/usr/bin/env bash

# Early exit test
# test whether any worker or coordinator exits before the
# task has completed (i.e., all output files have been finalized)
#
# This script assumes the following:
#   - in PATH: coordinator, worker
#   - in the current directory: early_exit.so
#   - in ENV: DATA_DIR (test data path)

echo '***' Starting early exit test.

exit 1

timeout -k 2s 180s coordinator "$DATA_DIR"/pg*txt &

# give the coordinator time to create the sockets.
sleep 1

# start multiple workers.
timeout -k 2s 180s worker early_exit.so &
timeout -k 2s 180s worker early_exit.so &
timeout -k 2s 180s worker early_exit.so &

# wait for any of the coordinator or workers to exit
# `jobs` ensures that any completed old processes from other tests
# are not waited upon
jobs &> /dev/null
wait -n

# a process has exited. this means that the output should be finalized
# otherwise, either a worker or the coordinator exited early
sort mr-out* | grep . > mr-wc-all-initial

# wait for remaining workers and coordinator to exit.
wait

# compare initial and final outputs
sort mr-out* | grep . > mr-wc-all-final
if cmp mr-wc-all-final mr-wc-all-initial
then
  echo '---' early exit test: PASS
  exit 0
else
  echo '---' output changed after first worker exited
  echo '---' early exit test: FAIL
  exit 1
fi
