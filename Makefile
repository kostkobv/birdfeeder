#########################
# PROJECT CONFIGURATION #
#########################
SHELL := /bin/bash

SOURCEDIR=.

VERSION=1.0.0

PROJECT_SOURCE=./

#####################
# TASKS DESCRIPTION #
#####################
.PHONY: all lint test

all: lint test

###########
# LINTING #
###########

# Run linter
lint:
	@gometalinter.v1 --exclude="(mocks|vendor)" --disable=gotype --fast --deadline=360s ./...

#########
# TESTS #
#########

# Run tests
test:
	@go test -cover -race $(shell glide novendor | grep -Ev '(mocks|config)')

# Tests for the CI
testci:
	bash run_tests.sh

# Build binaries
build:
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/birdfeeder-darwin-x86_64 main.go
	@GOOS=linux GOARCH=amd64 go build -o ./bin/birdfeeder-linux-x86_64 main.go
	@GOOS=windows GOARCH=amd64 go build -o ./bin/birdfeeder-windows-x86_64.exe main.go