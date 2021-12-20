# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

RUN mkdir "app"
ADD . /app

## We specify that we now wish to execute
## any further commands inside our /app
## directory
WORKDIR /app

## we run go build to compile the binary
## executable of our Go program
RUN go build -o main ./cmd/main.go

## Our start command which kicks off
## our newly created binary executable
CMD ["/app/main"]