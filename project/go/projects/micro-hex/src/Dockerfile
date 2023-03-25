# syntax=docker/dockerfile:1

FROM golang:{{ .Go.Version }}-alpine AS build

RUN apk add build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY / ./
RUN go test
RUN go build -o ./{{ .Name | kebab}}

FROM alpine

WORKDIR /go

COPY --from=build /app/{{ .Name | kebab}} ./{{ .Name | kebab}}
COPY *.json version ./

ENTRYPOINT ["/go/{{ .Name | kebab}}"]
