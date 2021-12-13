# syntax=docker/dockerfile:1

FROM golang:1.17-bullseye as base

WORKDIR /app

ENV GO111MODULE=on

# Comment this to run the tests without the Go race detector.
ENV RACE=-race


FROM base as build

COPY . ./

RUN go build $RACE ./cmd/coordinator

RUN go build $RACE ./cmd/worker

RUN go build $RACE ./test/mrsequential.go


FROM base as test

COPY ./test/testdata ./testdata

COPY --from=build /app/coordinator /usr/local/bin

COPY --from=build /app/worker /usr/local/bin

COPY --from=build /app/mrsequential /usr/local/bin
