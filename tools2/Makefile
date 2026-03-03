GO ?= go
BIN_DIR ?= bin
BIN_NAME ?= tools2
BIN_PATH := $(BIN_DIR)/$(BIN_NAME)

.PHONY: help setup generate tidy format lint build test run check

help:
	@echo "Available targets:"
	@echo "  make setup     - install templ/staticcheck tools"
	@echo "  make generate  - generate templ files"
	@echo "  make tidy      - run go mod tidy"
	@echo "  make format    - run go fmt ./..."
	@echo "  make lint      - run go vet + staticcheck"
	@echo "  make build     - build binary to $(BIN_PATH)"
	@echo "  make test      - run go test ./..."
	@echo "  make run       - start app on :8080"
	@echo "  make check     - generate + tidy + format + lint + build + test"

setup:
	$(GO) get -tool github.com/a-h/templ/cmd/templ@latest
	$(GO) get -tool honnef.co/go/tools/cmd/staticcheck@latest

generate:
	$(GO) tool templ generate ./...

tidy:
	$(GO) mod tidy

format:
	$(GO) fmt ./...

lint:
	$(GO) vet ./...
	$(GO) tool staticcheck ./...

build:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_PATH) .

test:
	$(GO) test ./...

run:
	$(GO) run .

check: generate tidy format lint build test
