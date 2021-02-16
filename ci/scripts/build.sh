#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-reporter-client
  make build
popd