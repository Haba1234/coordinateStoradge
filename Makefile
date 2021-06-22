BIN := "./bin/service"

build:
	go build -v -o $(BIN) ./cmd/service

run: build
	$(BIN) -port 8080

test:
	go test -v -count=100 -race -timeout=5m ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.39.0

lint: install-lint-deps
	golangci-lint run ./...

generate:
	go generate ./...

CLIENT_BIN := "./bin/client"

client-build:
	go build -v -o $(CLIENT_BIN) ./cmd/client

client-run: client-build
	$(CLIENT_BIN)

.PHONY: build run test lint generate client-run