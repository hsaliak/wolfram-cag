# Top-level Makefile for wolfram-cag

BINARY := wolfram-cag
CMD_PKG := ./cmd/wolfram-cag

.PHONY: all build test clean

all: build

build:
	go build -o $(BINARY) $(CMD_PKG)

test:
	go test ./...

clean:
	rm -f $(BINARY)
	go clean -testcache
