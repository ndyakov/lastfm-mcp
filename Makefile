.PHONY: build check test

build:
	go build -o bin/lastfm-mcp ./cmd/lastfm-mcp

test:
	go test -race -coverprofile=coverage.txt ./...

check:
	test -z "$$(gofmt -l .)"
	go vet ./...
	go test -race ./...
