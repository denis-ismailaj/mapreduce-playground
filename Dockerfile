# syntax=docker/dockerfile:1

FROM golang:1.17.3-bullseye

WORKDIR /app

COPY src/. ./

WORKDIR /app/main

RUN chmod +x test-mr.sh

RUN bash test-mr.sh
