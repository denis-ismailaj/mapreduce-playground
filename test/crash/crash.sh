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
sort out/mr-out-0 >mr-correct-crash.txt
rm -f out/mr-out*

(
  timeout -k 2s 180s coordinator "$DATA_DIR"/pg*txt
  touch mr-done
) &
sleep 1

# start multiple workers
timeout -k 2s 180s worker crash.so &

HOST='127.0.0.1'
PORT='1234'

function checkc() {
  nc -v -z $HOST $PORT &>/dev/null
}

(while checkc && [ ! -f mr-done ]; do
  timeout -k 2s 180s worker crash.so
  sleep 1
done) &

(while checkc && [ ! -f mr-done ]; do
  timeout -k 2s 180s worker crash.so
  sleep 1
done) &

while checkc && [ ! -f mr-done ]; do
  timeout -k 2s 180s worker crash.so
  sleep 1
done

wait

sort out/mr-out* | grep . >mr-crash-all
if cmp mr-crash-all mr-correct-crash.txt; then
  echo '---' crash test: PASS
  exit 0
else
  echo '---' crash output is not the same as mr-correct-crash.txt
  echo '---' crash test: FAIL
  exit 1
fi
