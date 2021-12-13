#!/usr/bin/env bash

# Indexer test
#
# This script assumes the following:
#   - in PATH: coordinator, worker, mrsequential
#   - in the current directory: indexer.so
#   - in ENV: DATA_DIR (test data path)

# generate the correct output
mrsequential indexer.so "$DATA_DIR"/pg*txt || exit 1
sort out/mr-out-0 >mr-correct-indexer.txt
rm -f out/mr-out*

echo '***' Starting indexer test.

timeout -k 2s 180s coordinator "$DATA_DIR"/pg*txt &
sleep 1

# start multiple workers
timeout -k 2s 180s worker indexer.so &
timeout -k 2s 180s worker indexer.so

sort out/mr-out* | grep . >mr-indexer-all
if cmp mr-indexer-all mr-correct-indexer.txt; then
  echo '---' indexer test: PASS
  exit 0
else
  echo '---' indexer output is not the same as mr-correct-indexer.txt
  echo '---' indexer test: FAIL
  exit 1
fi
