gexport GOBIN := $(shell pwd)/bin
export PATH := $(GOBIN):$(PATH)

build: 
	go build -o bin/svc ./cmd/server

tools: 
	go install github.com/matryer/moq@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

.PHONY: lint
lint:
	golangci-lint run

test: 
	go test -count=1 ./...

test_integration:
	INTEGRATION_TEST=on DB_DSN=./tmp/test_db.sqlite3 go test -count=1 -v -run Integration ./...

generate: tools
	go generate ./...
