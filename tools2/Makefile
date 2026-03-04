GO ?= go
BIN_DIR ?= bin
BIN_NAME ?= tools2
BIN_PATH := $(BIN_DIR)/$(BIN_NAME)
APP_DIR := ./app

.PHONY: help setup generate tidy format lint build test run dev check

help:
	@echo "Available targets:"
	@echo "  make setup     - install templ/staticcheck tools"
	@echo "  make generate  - generate templ files"
	@echo "  make tidy      - run go mod tidy"
	@echo "  make format    - run go fmt ./..."
	@echo "  make lint      - run go vet + staticcheck"
	@echo "  make build     - build binary to $(BIN_PATH)"
	@echo "  make test      - run go test ./..."
	@echo "  make run       - generate templ and start app on :8080"
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

test:
	$(GO) test ./...

run: generate
	echo http://localhost:8080
	$(GO) run $(APP_DIR)

dev: run

check: generate tidy format lint build test
