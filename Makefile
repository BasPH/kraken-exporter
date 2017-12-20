SHELL := /bin/bash
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all test fmt build vet lint

all: test fmt vet lint build

test:
	@echo do something

fmt:
	@gofmt -s -l -w $(SRC)

build:
	@echo do something

vet:
	@echo do something

lint:
	@echo do something