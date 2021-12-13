#!/usr/bin/env bash

# Map parallelism test
#
# This script assumes the following:
#   - in PATH: coordinator, worker
#   - in the current directory: mtiming.so
#   - in ENV: DATA_DIR (test data path)

echo '***' Starting map parallelism test.

timeout -k 2s 180s coordinator "$DATA_DIR"/pg*txt &
sleep 1

timeout -k 2s 180s worker mtiming.so &
timeout -k 2s 180s worker mtiming.so

NT=$(cat out/mr-out* | grep -c '^times-' | sed 's/ //g')
if [ "$NT" != "2" ]
then
  echo '---' saw "$NT" workers rather than 2
  echo '---' map parallelism test: FAIL
  exit 1
fi

if cat out/mr-out* | grep '^parallel.* 2' > /dev/null
then
  echo '---' map parallelism test: PASS
  exit 0
else
  echo '---' map workers did not run in parallel
  echo '---' map parallelism test: FAIL
  exit 1
fi
