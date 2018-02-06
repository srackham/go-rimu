# Simple makefile to log workflow.

.ONESHELL:
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

push:
	git push -u origin master

# Run fuzz test.
fuzz:
	go-fuzz -bin=./rimu/rimu-fuzz.zip -workdir=./rimu/fuzz-workdir

# Build fuzz executables.
# If the fuzz work directory does not exist it is created and the corpus is primed.
# Requires go-fuzz package (https://github.com/dvyukov/go-fuzz):
fuzz-build:
	@set -eu
	cd rimu
	if [ ! -d fuzz-workdir/corpus ]; then
		echo Creating workdir...
		mkdir -p fuzz-workdir/corpus
		unzip -q testdata/fuzz-samples.zip -d fuzz-workdir/corpus
	fi
	echo Building executables...
	go-fuzz-build github.com/srackham/go-rimu/rimu

# List fuzz crash inputs.
fuzz-crashes:
	@set -eu
	for f in rimu/fuzz-workdir/crashers/*.quoted; do
		echo $$f
		cat $$f
	done