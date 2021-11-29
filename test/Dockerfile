# syntax=docker/dockerfile:1

FROM golang:1.17.3-bullseye

WORKDIR /app

COPY . ./

WORKDIR /app/test/

RUN chmod +x test-mr.sh

RUN bash test-mr.sh
