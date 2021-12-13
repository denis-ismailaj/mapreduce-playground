#!/usr/bin/env bash

# Reduce parallelism test
#
# This script assumes the following:
#   - in PATH: coordinator, worker
#   - in the current directory: rtiming.so
#   - in ENV: DATA_DIR (test data path)

echo '***' Starting reduce parallelism test.

timeout -k 2s 180s coordinator "$DATA_DIR"/pg*txt &
sleep 1

timeout -k 2s 180s worker rtiming.so &
timeout -k 2s 180s worker rtiming.so

NT=$(cat out/mr-out* | grep -c '^[a-z] 2' | sed 's/ //g')
if [ "$NT" -lt "2" ]
then
  echo '---' too few parallel reduces.
  echo '---' reduce parallelism test: FAIL
  exit 1
else
  echo '---' reduce parallelism test: PASS
  exit 0
fi
