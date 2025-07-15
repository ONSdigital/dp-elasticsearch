SHELL=bash

all: audit lint test build
.PHONY: all

test:
	go test -race -cover ./...
.PHONY: test

audit:
	go list -json -m all | nancy sleuth
.PHONY: audit

build:
	go build ./...
.PHONY: build

lint:
	golangci-lint run ./...
.PHONY: lint
