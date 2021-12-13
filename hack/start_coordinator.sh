#!/usr/bin/env bash

# This is a workaround when using the project locally with Docker Compose
# The coordinator requires the input files to be run, but the volume with the input files
# is only mounted after the entrypoint is run.
#
# This script waits until the volume is ready, then runs the coordinator

until cd /app/inputs
do
    echo "Waiting for bind"
done

/coordinator /app/inputs/*
