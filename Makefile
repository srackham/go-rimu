# Simple makefile to log workflow.

# Set defaults (see http://clarkgrubb.com/makefile-style-guide#prologue)
MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:
.ONESHELL:

GOFLAGS ?= $(GOFLAGS:)

.PHONY: all
all: test install

./rimugo/bindata.go: ./rimugo/resources/*
	cd ./rimugo && go-bindata ./resources

.PHONY: bindata
bindata: ./rimugo/bindata.go

.PHONY: install
install: bindata
	go install $(GOFLAGS) ./...

.PHONY: build
build: bindata
	go build $(GOFLAGS) ./...

.PHONY: test
test: bindata
	go test $(GOFLAGS) ./...

.PHONY: clean
clean:
	go clean $(GOFLAGS) -i ./...

.PHONY: push
push:
	git push -u --tags origin master

# Run fuzz test.
.PHONY: fuzz
fuzz:
	go-fuzz -bin=./rimu/rimu-fuzz.zip -workdir=./rimu/fuzz-workdir

# Build fuzz executables.
# If the fuzz work directory does not exist it is created and the corpus is primed.
# Requires go-fuzz package (https://github.com/dvyukov/go-fuzz):
.PHONY: fuzz-build
fuzz-build:
	@cd rimu
	if [ ! -d fuzz-workdir/corpus ]; then
		echo Initializing corpus...
		mkdir -p fuzz-workdir/corpus
		unzip -q testdata/fuzz-samples.zip -d fuzz-workdir/corpus
	fi
	echo Building executables...
	go-fuzz-build github.com/srackham/go-rimu/rimu

# List fuzz crash inputs.
.PHONY: fuzz-crashes
fuzz-crashes:
	@for f in rimu/fuzz-workdir/crashers/*.quoted; do
		echo $$f
		cat $$f
	done