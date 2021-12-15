# syntax=docker/dockerfile:1

FROM golang:1.17-bullseye as build

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /coordinator ./cmd/coordinator
RUN go build -o /worker ./cmd/worker

RUN chmod +x ./hack/start_coordinator.sh


FROM ubuntu

WORKDIR /

COPY --from=build /coordinator .
COPY --from=build /worker .

COPY --from=build /app/hack/start_coordinator.sh .
