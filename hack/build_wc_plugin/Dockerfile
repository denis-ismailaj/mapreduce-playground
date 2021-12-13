# syntax=docker/dockerfile:1

FROM golang:1.17-bullseye as build

WORKDIR /app

ENV GO111MODULE=on

COPY . .

ENTRYPOINT go build -buildmode=plugin -o out/plugin.so ./test/wc/wc.go
