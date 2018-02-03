# Simple makefile to log workflow.

.PHONY: all bindata install build test clean fuzz

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

fuzz:
	# Before you can run fuzz you need to install the go-fuzz package (https://github.com/dvyukov/go-fuzz):
	#
	#    go get github.com/dvyukov/go-fuzz/go-fuzz
	#    go get github.com/dvyukov/go-fuzz/go-fuzz-build
	#
	# Then generate the fuzz execuatables and and working directories:
	#
	#    cd rimu
	#    go-fuzz-build github.com/srackham/go-rimu/rimu
	#    mkdir fuzz-workdir
	#
	go-fuzz -bin=./rimu/rimu-fuzz.zip -workdir=./rimu/fuzz-workdir