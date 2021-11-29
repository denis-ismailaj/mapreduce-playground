#!/usr/bin/env bash

#
# basic map-reduce test
#

#RACE=

# comment this to run the tests without the Go race detector.
RACE=-race

MR_APPS_DIR=mrapps
COORDINATOR=../cmd/coordinator/main.go
WORKER=../cmd/worker/main.go
SEQUENTIAL=./mrsequential.go

# run the test in a fresh sub-directory.
rm -rf mr-tmp
mkdir mr-tmp || exit 1
cd mr-tmp || exit 1
rm -f mr-*

# make sure software is freshly built.
(cd $MR_APPS_DIR && go build $RACE -buildmode=plugin wc.go) || exit 1
(cd $MR_APPS_DIR && go build $RACE -buildmode=plugin indexer.go) || exit 1
(cd $MR_APPS_DIR && go build $RACE -buildmode=plugin mtiming.go) || exit 1
(cd $MR_APPS_DIR && go build $RACE -buildmode=plugin rtiming.go) || exit 1
(cd $MR_APPS_DIR && go build $RACE -buildmode=plugin jobcount.go) || exit 1
(cd $MR_APPS_DIR && go build $RACE -buildmode=plugin early_exit.go) || exit 1
(cd $MR_APPS_DIR && go build $RACE -buildmode=plugin crash.go) || exit 1
(cd $MR_APPS_DIR && go build $RACE -buildmode=plugin nocrash.go) || exit 1
(cd .. && go build $RACE $COORDINATOR) || exit 1
(cd .. && go build $RACE $WORKER) || exit 1
(cd .. && go build $RACE $SEQUENTIAL) || exit 1

failed_any=0

#########################################################
# first word-count

# generate the correct output
../mrsequential $MR_APPS_DIR/wc.so ../pg*txt || exit 1
sort mr-out-0 > mr-correct-wc.txt
rm -f mr-out*

echo '***' Starting wc test.

timeout -k 2s 180s ../mrcoordinator ../pg*txt &
pid=$!

# give the coordinator time to create the sockets.
sleep 1

# start multiple workers.
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/wc.so &
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/wc.so &
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/wc.so &

# wait for the coordinator to exit.
wait $pid

# since workers are required to exit when a job is completely finished,
# and not before, that means the job has finished.
sort mr-out* | grep . > mr-wc-all
if cmp mr-wc-all mr-correct-wc.txt
then
  echo '---' wc test: PASS
else
  echo '---' wc output is not the same as mr-correct-wc.txt
  echo '---' wc test: FAIL
  failed_any=1
fi

# wait for remaining workers and coordinator to exit.
wait

#########################################################
# now indexer
rm -f mr-*

# generate the correct output
../mrsequential $MR_APPS_DIR/indexer.so ../pg*txt || exit 1
sort mr-out-0 > mr-correct-indexer.txt
rm -f mr-out*

echo '***' Starting indexer test.

timeout -k 2s 180s ../mrcoordinator ../pg*txt &
sleep 1

# start multiple workers
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/indexer.so &
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/indexer.so

sort mr-out* | grep . > mr-indexer-all
if cmp mr-indexer-all mr-correct-indexer.txt
then
  echo '---' indexer test: PASS
else
  echo '---' indexer output is not the same as mr-correct-indexer.txt
  echo '---' indexer test: FAIL
  failed_any=1
fi

wait

#########################################################
echo '***' Starting map parallelism test.

rm -f mr-*

timeout -k 2s 180s ../mrcoordinator ../pg*txt &
sleep 1

timeout -k 2s 180s ../mrworker $MR_APPS_DIR/mtiming.so &
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/mtiming.so

NT=$(cat mr-out* | grep -c '^times-' | sed 's/ //g')
if [ "$NT" != "2" ]
then
  echo '---' saw "$NT" workers rather than 2
  echo '---' map parallelism test: FAIL
  failed_any=1
fi

if cat mr-out* | grep '^parallel.* 2' > /dev/null
then
  echo '---' map parallelism test: PASS
else
  echo '---' map workers did not run in parallel
  echo '---' map parallelism test: FAIL
  failed_any=1
fi

wait


#########################################################
echo '***' Starting reduce parallelism test.

rm -f mr-*

timeout -k 2s 180s ../mrcoordinator ../pg*txt &
sleep 1

timeout -k 2s 180s ../mrworker $MR_APPS_DIR/rtiming.so &
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/rtiming.so

NT=$(cat mr-out* | grep -c '^[a-z] 2' | sed 's/ //g')
if [ "$NT" -lt "2" ]
then
  echo '---' too few parallel reduces.
  echo '---' reduce parallelism test: FAIL
  failed_any=1
else
  echo '---' reduce parallelism test: PASS
fi

wait

#########################################################
echo '***' Starting job count test.

rm -f mr-*

timeout -k 2s 180s ../mrcoordinator ../pg*txt &
sleep 1

timeout -k 2s 180s ../mrworker $MR_APPS_DIR/jobcount.so &
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/jobcount.so
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/jobcount.so &
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/jobcount.so

NT=$(cat mr-out* | awk '{print $2}')
if [ "$NT" -ne "8" ]
then
  echo '---' map jobs ran incorrect number of times "($NT != 8)"
  echo '---' job count test: FAIL
  failed_any=1
else
  echo '---' job count test: PASS
fi

wait

#########################################################
# test whether any worker or coordinator exits before the
# task has completed (i.e., all output files have been finalized)
rm -f mr-*

echo '***' Starting early exit test.

timeout -k 2s 180s ../mrcoordinator ../pg*txt &

# give the coordinator time to create the sockets.
sleep 1

# start multiple workers.
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/early_exit.so &
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/early_exit.so &
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/early_exit.so &

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
else
  echo '---' output changed after first worker exited
  echo '---' early exit test: FAIL
  failed_any=1
fi
rm -f mr-*

#########################################################
echo '***' Starting crash test.

# generate the correct output
../mrsequential $MR_APPS_DIR/nocrash.so ../pg*txt || exit 1
sort mr-out-0 > mr-correct-crash.txt
rm -f mr-out*

rm -f mr-done
(timeout -k 2s 180s ../mrcoordinator ../pg*txt ; touch mr-done ) &
sleep 1

# start multiple workers
timeout -k 2s 180s ../mrworker $MR_APPS_DIR/crash.so &

# mimic rpc.go's coordinatorSock()
SOCK_NAME=/var/tmp/824-mr-$(id -u)

( while [ -e "$SOCK_NAME" ] && [ ! -f mr-done ]
  do
    timeout -k 2s 180s ../mrworker $MR_APPS_DIR/crash.so
    sleep 1
  done ) &

( while [ -e "$SOCK_NAME" ] && [ ! -f mr-done ]
  do
    timeout -k 2s 180s ../mrworker $MR_APPS_DIR/crash.so
    sleep 1
  done ) &

while [ -e "$SOCK_NAME" ] && [ ! -f mr-done ]
do
  timeout -k 2s 180s ../mrworker $MR_APPS_DIR/crash.so
  sleep 1
done

wait

rm "$SOCK_NAME"
sort mr-out* | grep . > mr-crash-all
if cmp mr-crash-all mr-correct-crash.txt
then
  echo '---' crash test: PASS
else
  echo '---' crash output is not the same as mr-correct-crash.txt
  echo '---' crash test: FAIL
  failed_any=1
fi

#########################################################
if [ $failed_any -eq 0 ]; then
    echo '***' PASSED ALL TESTS
else
    echo '***' FAILED SOME TESTS
    exit 1
fi
