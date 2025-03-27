SHELL=bash

all: audit test build lint
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
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5
	golangci-lint run ./...
.PHONY: lint
