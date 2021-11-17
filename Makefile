SHELL=bash

all: audit test build lint
.PHONY: all

.PHONY: fmt
fmt:
	go fmt ./...

test:
	go test -race -cover ./...
.PHONY: test

audit:
	go list -json -m all | nancy sleuth
.PHONY: audit

build: fmt
	go build ./...
.PHONY: build

lint:
	golangci-lint --deadline=10m --fast --enable=gocritic --enable=gofmt --enable=gocyclo --enable=bodyclose --enable=gocognit run
	golangci-lint --fast --tests=false --enable=gosec run
.PHONY: lint
