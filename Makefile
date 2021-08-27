# Makefile to log workflow.

# Set defaults (see http://clarkgrubb.com/makefile-style-guide#prologue)
MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:
.ONESHELL:
.SILENT:

BUILDFLAGS := -ldflags "-s -w"
VERS := $$(sed -ne 's/^const VERSION = "\([0-9a-z.]*\)"/\1/p' rimugo/rimugo.go)

.PHONY: all
all: test install

./rimugo/bindata.go: ./rimugo/resources/*
	cd ./rimugo && go-bindata ./resources

.PHONY: bindata
bindata: ./rimugo/bindata.go

.PHONY: install
install: bindata
	go install $(BUILDFLAGS) ./...

.PHONY: build
build: bindata
	go build $(BUILDFLAGS) ./...

.PHONY: test
test: build
	go test ./...

.PHONY: clean
clean:
	go clean -i ./...
	git gc --prune=now

.PHONY: tag
tag: test
	tag=v$(VERS)
	echo tag: $$tag
	git tag -a -m "$$tag" "$$tag"

.PHONY: push
push:
	git push -u --tags origin master

# Run fuzz test.
.PHONY: fuzz
fuzz: fuzz-corpus fuzz-build
	go-fuzz -bin=./rimu/rimu-fuzz.zip -workdir=./rimu/fuzz-workdir

# Build fuzz executables.
.PHONY: fuzz-build
fuzz-build: ./rimu/rimu-fuzz.zip

./rimu/rimu-fuzz.zip: $(shell find . -name '*.go')
	@echo Building executables...
	go-fuzz-build -o ./rimu/rimu-fuzz.zip github.com/srackham/go-rimu/rimu

# If the fuzz corpus directory does not exist it is created and the corpus is primed.
.PHONY: fuzz-corpus
fuzz-corpus:
	@cd rimu
	if [ ! -d fuzz-workdir/corpus ]; then
		echo Initializing corpus...
		mkdir -p fuzz-workdir/corpus
		unzip -q testdata/fuzz-samples.zip -d fuzz-workdir/corpus
	fi

# List fuzz crashes.
.PHONY: fuzz-crashes
fuzz-crashes:
	@for f in rimu/fuzz-workdir/crashers/*.quoted; do
		echo $$f
		cat $$f
	done

.PHONY: benchmark
benchmark:
	cd rimu
	go test -bench .

.PHONY: profile
profile:
	cd rimu
	go test -cpuprofile cpu.out -bench .    # Create benchmark executable and profile.
	go tool pprof -text rimu.test cpu.out   # View top profile entries.