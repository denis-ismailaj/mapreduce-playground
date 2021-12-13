#!/usr/bin/env bash

# Job count test
#
# This script assumes the following:
#   - in PATH: coordinator, worker
#   - in the current directory: jobcount.so
#   - in ENV: DATA_DIR (test data path)

echo '***' Starting job count test.

timeout -k 2s 180s coordinator "$DATA_DIR"/pg*txt &
sleep 1

timeout -k 2s 180s worker jobcount.so &
timeout -k 2s 180s worker jobcount.so
timeout -k 2s 180s worker jobcount.so &
timeout -k 2s 180s worker jobcount.so

NT=$(cat out/mr-out* | awk '{print $2}')
if [ "$NT" -ne "8" ]
then
  echo '---' map jobs ran incorrect number of times "($NT != 8)"
  echo '---' job count test: FAIL
  exit 1
else
  echo '---' job count test: PASS
  exit 0
fi
