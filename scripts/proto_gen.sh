#!/bin/bash

protoc -I=proto/ --go_out=pkg/ proto/liveness.proto
protoc -I=proto/ --go_out=pkg/ proto/pbft.proto
protoc -I=proto/ --go_out=pkg/ proto/controller.proto
protoc -I=proto/ --go-grpc_out=pkg/ proto/liveness.proto
protoc -I=proto/ --go-grpc_out=pkg/ proto/pbft.proto
protoc -I=proto/ --go-grpc_out=pkg/ proto/controller.proto
