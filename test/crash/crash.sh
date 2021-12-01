#!/usr/bin/env bash

# Crash test
#
# This script assumes the following:
#   - in PATH: coordinator, worker, mrsequential
#   - in the current directory: crash.so, nocrash.so
#   - in ENV: DATA_DIR (test data path)

echo '***' Starting crash test.

# generate the correct output
mrsequential nocrash.so "$DATA_DIR"/pg*txt || exit 1
sort mr-out-0 > mr-correct-crash.txt
rm -f mr-out*

rm -f mr-done
(timeout -k 2s 180s coordinator "$DATA_DIR"/pg*txt ; touch mr-done ) &
sleep 1

# start multiple workers
timeout -k 2s 180s worker crash.so &

# mimic rpc.go's coordinatorSock()
SOCK_NAME=/var/tmp/824-mr-$(id -u)

( while [ -e "$SOCK_NAME" ] && [ ! -f mr-done ]
  do
    timeout -k 2s 180s worker crash.so
    sleep 1
  done ) &

( while [ -e "$SOCK_NAME" ] && [ ! -f mr-done ]
  do
    timeout -k 2s 180s worker crash.so
    sleep 1
  done ) &

while [ -e "$SOCK_NAME" ] && [ ! -f mr-done ]
do
  timeout -k 2s 180s worker crash.so
  sleep 1
done

wait

rm "$SOCK_NAME"
sort mr-out* | grep . > mr-crash-all
if cmp mr-crash-all mr-correct-crash.txt
then
  echo '---' crash test: PASS
  exit 0
else
  echo '---' crash output is not the same as mr-correct-crash.txt
  echo '---' crash test: FAIL
  exit 1
fi
