SHELL := /bin/bash
GOFILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all test fmt build vet lint

all: test fmt vet lint build

test:
	@echo do something

fmt:
	@gofmt -s -l -w $(GOFILES)

build:
	@echo do something

vet:
	@go tool vet $(GOFILES)

lint:
	@for f in $(GOFILES); do golint $d; done