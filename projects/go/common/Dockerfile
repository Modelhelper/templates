# syntax=docker/dockerfile:1

FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY / ./
COPY version ./

RUN go build -o ./logging-proxy

FROM alpine

RUN apk add curl

WORKDIR /
COPY --from=build /app/logging-proxy /logging-proxy
COPY *.json /
COPY version /

EXPOSE 8686

ENTRYPOINT ["/logging-proxy"]