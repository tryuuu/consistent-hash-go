SHELL := /bin/bash

.PHONY: run test

run:
	go run ./cmd

test:
	go test ./...

