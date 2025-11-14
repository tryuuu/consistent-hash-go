SHELL := /bin/bash

.PHONY: run test benchmark

run:
	go run ./cmd

test:
	go test ./...

benchmark:
	go run ./cmd/benchmark

