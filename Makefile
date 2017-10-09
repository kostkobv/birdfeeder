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
	@gometalinter.v1 -e vendor --deadline=60s --vendor $(PROJECT_SOURCE)

#########
# TESTS #
#########

test:
	@go test $(shell go list ./... | grep -v vendor) -cover
