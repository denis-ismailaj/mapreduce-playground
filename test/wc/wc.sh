#!/usr/bin/env bash

# Word-count test
#
# This script assumes the following:
#   - in PATH: coordinator, worker, mrsequential
#   - in the current directory: wc.so
#   - in ENV: DATA_DIR (test data path)

# generate the correct output
mrsequential wc.so "$DATA_DIR"/pg*txt || exit 1
sort out/mr-out-0 >mr-correct-wc.txt
rm -f out/mr-out*

echo '***' Starting wc test.

timeout -k 2s 180s coordinator "$DATA_DIR"/pg*txt &
pid=$!

# give the coordinator time to create the sockets.
sleep 1

# start multiple workers.
timeout -k 2s 180s worker wc.so &
timeout -k 2s 180s worker wc.so &
timeout -k 2s 180s worker wc.so &

# wait for the coordinator to exit.
wait $pid

# since workers are required to exit when a job is completely finished,
# and not before, that means the job has finished.
sort out/mr-out* | grep . >mr-wc-all
if cmp mr-wc-all mr-correct-wc.txt; then
  echo '---' wc test: PASS
  exit 0
else
  echo '---' wc output is not the same as mr-correct-wc.txt
  echo '---' wc test: FAIL
  exit 1
fi
