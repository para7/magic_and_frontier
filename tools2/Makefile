GO ?= go
BIN_DIR ?= bin
BIN_NAME ?= tools2
BIN_PATH := $(BIN_DIR)/$(BIN_NAME)
APP_DIR := ./app
CLI_BIN_NAME ?= tools2-cli
CLI_BIN_PATH := $(BIN_DIR)/$(CLI_BIN_NAME)
CLI_DIR := ./app/cmd/tools2-cli

.PHONY: help setup generate tidy format lint build build-cli test run run-cli dev check

help:
	@echo "Available targets:"
	@echo "  make setup     - install templ/staticcheck tools"
	@echo "  make generate  - generate templ files"
	@echo "  make tidy      - run go mod tidy"
	@echo "  make format    - run go fmt ./..."
	@echo "  make lint      - run go vet + staticcheck"
	@echo "  make build     - build binary to $(BIN_PATH)"
	@echo "  make build-cli - build CLI binary to $(CLI_BIN_PATH)"
	@echo "  make test      - run go test ./..."
	@echo "  make run       - generate templ and start app on :8080"
	@echo "  make run-cli   - run CLI binary with ARGS='...'"
	@echo "  make dev       - alias of make run"
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
	$(GO) build -o $(BIN_PATH) $(APP_DIR)

build-cli:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(CLI_BIN_PATH) $(CLI_DIR)

test:
	$(GO) test ./...

run: generate
	echo http://localhost:8080
	$(GO) run $(APP_DIR)

run-cli:
	$(GO) run $(CLI_DIR) $(ARGS)

dev: run

check: generate tidy format lint build test
