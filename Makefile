export GOBIN := $(shell pwd)/bin
export PATH := $(GOBIN):$(PATH)

tools: 
	go install github.com/matryer/moq@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

.PHONY: lint
lint:
	golangci-lint run

test: 
	go test -count=1 ./...

test_integration:
	INTEGRATION_TEST=on go test -count=1 -v -run Integration ./...

generate: tools
	go generate ./...
