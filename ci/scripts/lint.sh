#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-elasticsearch
# Install golangci-lint
  go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1
  make lint
popd
