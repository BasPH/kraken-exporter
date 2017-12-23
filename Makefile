SHELL := /bin/bash
GOFILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GIT_COMMIT=$(shell git rev-parse --short HEAD)
IMAGE_NAME := basph/kraken-exporter

.PHONY: all fmt vet lint build

all: fmt vet lint build

fmt:
	@gofmt -s -l -w $(GOFILES)

vet:
	@go tool vet $(GOFILES)

lint:
	@for f in $(GOFILES); do golint $d; done

build:
	@echo "Building Kraken-exporter git tag ${GIT_COMMIT}"
	docker build -t $(IMAGE_NAME) -f Dockerfile .
	docker tag $(IMAGE_NAME) $(IMAGE_NAME):latest
	docker tag $(IMAGE_NAME) $(IMAGE_NAME):${GIT_COMMIT}
