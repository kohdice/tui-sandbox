BIN := ./bin/sandbox
CMD := main.go
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS := "-s -w -X main.revision=$(CURRENT_REVISION)"

.PHONY: all
all: clean tidy build

.PHONY: build
build:
	CGO_ENABLED=0 go build -trimpath -ldflags=$(BUILD_LDFLAGS) -o $(BIN) $(CMD)

.PHONY: clean
clean:
	go clean

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	goimports -local github.com/kohdice/tui-sandbix -w .
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: coverage
coverage:
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
