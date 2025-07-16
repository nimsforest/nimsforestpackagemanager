# Main Project Makefile

include MAKEFILE.nimsforestpm
# This project uses nimsforestwebstack for web presence
-include ../nimsforestwebstack-workspace/main/Makefile.nimsforestwebstack

# Test targets
.PHONY: test test-integration test-all build-and-test

# Build the main binary
build:
	go build -o nimsforestpm ./cmd

# Build example tool for integration testing
build-example:
	go build -o bin/example-tool ./integration/example-tool

# Run unit tests only
test:
	go test ./...

# Run integration tests only 
test-integration:
	go test -tags=integration ./...

# Build and run integration tests (proper full cycle)
build-and-test: build build-example
	go test -v ./integration

# Run both unit and integration tests
test-all:
	go test ./...
	go test -tags=integration ./...
