FROM golang:1.12.4-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates git
WORKDIR /go/src/github.com/fmagege/tut-meower

ENV GO111MODULE=on

COPY go.mod go.mod
COPY go.sum go.sum

COPY db db
COPY event event
COPY meow-service meow-service
COPY pusher-service pusher-service
COPY query-service query-service
COPY schema schema
COPY search search
COPY util util

RUN go install ./...

FROM alpine:3.9
WORKDIR /usr/bin
COPY --from=build /go/bin .
