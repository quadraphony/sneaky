.PHONY: build test clean run validate

APP_NAME=sneakycli

build:
	go build -o bin/$(APP_NAME) ./cmd/sneakycli

test:
	go test ./...

run:
	go run ./cmd/sneakycli

validate:
	go test ./... && go build ./...

clean:
	rm -rf bin
