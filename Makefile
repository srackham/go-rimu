# Simple makefile to log workflow.

.PHONY: all bindata install build test clean fuzz fuzz-build fuzz-crashes

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
	go-fuzz -bin=./rimu/rimu-fuzz.zip -workdir=./rimu/fuzz-workdir

fuzz-build:
	# Before you can run fuzz you need to install the go-fuzz package (https://github.com/dvyukov/go-fuzz):
	#
	#    go get github.com/dvyukov/go-fuzz/go-fuzz
	#    go get github.com/dvyukov/go-fuzz/go-fuzz-build
	#
	# Create fuzz work directories; prime the corpus; build fuzz executables.
	cd rimu && \
	if [ -d fuzz-workdir ]; then rm -rf fuzz-workdir.OLD; mv fuzz-workdir fuzz-workdir.OLD; fi && \
	mkdir -p fuzz-workdir/corpus && \
	unzip -q testdata/fuzz-samples.zip -d fuzz-workdir/corpus && \
	go-fuzz-build github.com/srackham/go-rimu/rimu

fuzz-crashes:
	# List fuzz crash inputs.
	cat rimu/fuzz-workdir/crashers/*.quoted