.PHONY: all test fmt

all: vet test

vet:
	go vet ./...

fmt:
	find . -name '*.go' -type f -exec gofmt -w {} \;

test:
	go test ./...
