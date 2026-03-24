GO ?= go
BIN_DIR ?= bin
BIN_NAME ?= mce
BIN_PATH := $(BIN_DIR)/$(BIN_NAME)
APP_DIR := ./app
LINT_CACHE_DIR ?= /tmp/maf-command-editor-cache
GO_CACHE_DIR ?= /tmp/maf-command-editor-go-cache

.PHONY: help setup generate tidy format lint build test run run-cmd dev check

help:
	@echo "Available targets:"
	@echo "  make setup     - install templ/staticcheck tools"
	@echo "  make generate  - generate templ files"
	@echo "  make tidy      - run go mod tidy"
	@echo "  make format    - run go fmt ./..."
	@echo "  make lint      - run go vet + staticcheck"
	@echo "  make build     - build binary to $(BIN_PATH)"
	@echo "  make test      - run go test ./..."
	@echo "  make run       - generate templ and start editor on :8787"
	@echo "  make run-cmd   - run command with ARGS='validate|export'"
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
	mkdir -p $(GO_CACHE_DIR)
	mkdir -p $(LINT_CACHE_DIR)
	GOCACHE=$(GO_CACHE_DIR) $(GO) vet ./...
	GOCACHE=$(GO_CACHE_DIR) XDG_CACHE_HOME=$(LINT_CACHE_DIR) $(GO) tool staticcheck ./...

build:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_PATH) $(APP_DIR)

test:
	$(GO) test ./...

run: generate
	echo http://localhost:8787
	$(GO) run $(APP_DIR) editor

run-cmd:
	$(GO) run $(APP_DIR) $(ARGS)

dev: run

check: generate tidy format lint build test

PORT=8787
kill-port:
	@pids=$$(lsof -t -i :$(PORT) 2>/dev/null); \
	if [ -z "$$pids" ]; then \
		echo "No process is using port $(PORT)"; \
	else \
		echo "Killing PID(s): $$pids on port $(PORT)"; \
		kill $$pids; \
	fi
