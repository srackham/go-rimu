# Simple makefile to log workflow.
.PHONY: all bindata install build test clean

GOFLAGS ?= $(GOFLAGS:)

all: test install

./rimucgo/bindata.go: ./rimucgo/resources/*
	cd ./rimucgo && go-bindata ./resources

bindata: ./rimucgo/bindata.go

install: bindata
	go install $(GOFLAGS) ./...

build: bindata
	go build $(GOFLAGS) ./...

test: bindata
	go test $(GOFLAGS) ./...

clean:
	go clean $(GOFLAGS) -i ./...