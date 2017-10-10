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

lint:
	@gometalinter.v1 --exclude=vendor --disable=gotype --fast --deadline=360s ./...

#########
# TESTS #
#########

test:
	@go test
