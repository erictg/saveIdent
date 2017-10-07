#!/usr/bin/env bash

echo "this should be run from CD perspective of project top level directory!"

echo remember to set your GOPATH to:

echo $(pwd)

export GOPATH=$(pwd)

go get -u github.com/kardianos/govendor