#!/bin/bash

# Install protoc
# curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v27.2/protoc-27.2-linux-x86_64.zip

# Install plugins
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# export PATH="$PATH:$(go env GOPATH)/bin"

# Install mock
# go install github.com/golang/mock/mockgen@latest
# export PATH="$PATH:$(go env GOPATH)/bin"

# Build proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./drive/rpc/rpc.proto

# Build mock
# TODO
