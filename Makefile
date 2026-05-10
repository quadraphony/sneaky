# Makefile for sneaky core

.PHONY: all build test fmt vet lint docker-build bench
all: build

build:
	go build ./...

test:
	go test ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

lint:
	golangci-lint run ./...

docker-build:
	docker build -t quadraphony/sneaky:dev .

bench:
	go test -bench=. ./...
