# Simple makefile to log workflow.
.PHONY: all bindata build test clean

GOFLAGS ?= $(GOFLAGS:)

all: test

./rimucgo/bindata.go: ./rimucgo/resources/*
	cd ./rimucgo && go-bindata ./resources

bindata: ./rimucgo/bindata.go

build: bindata
	go build $(GOFLAGS) ./...

test: bindata
	go test $(GOFLAGS) ./...

clean:
	go clean $(GOFLAGS) -i ./...