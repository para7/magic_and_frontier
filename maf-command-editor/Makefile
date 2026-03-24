.PHONY: help generate tidy format lint build test run run-cmd dev check kill-port

help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: run/web ## generate templ and start editor on :8787

run/editor:
	echo http://localhost:8787
	go run ./app editor

run/validate:
	go run ./app validate

run/export:
	go run ./app export


check: generate tidy format lint build test ## run all checks

generate: ## generate templ files
	go tool templ generate ./...

tidy: ## run go mod tidy
	go mod tidy

format: ## run go fmt ./...
	go fmt ./...

lint: ## run go vet + staticcheck
	mkdir -p /tmp/maf-command-editor-go-cache
	mkdir -p /tmp/maf-command-editor-cache
	GOCACHE=/tmp/maf-command-editor-go-cache go vet ./...
	GOCACHE=/tmp/maf-command-editor-go-cache XDG_CACHE_HOME=/tmp/maf-command-editor-cache go tool staticcheck ./...

build: ## build binary to bin/mce
	mkdir -p bin
	go build -o bin/mce ./app

test: ## run go test ./...
	go test ./...



run-cmd: ## run command with ARGS='validate|export'
	go run ./app $(ARGS)

dev: run ## alias of make run

kill-port: ## kill process using port 8787
	@pids=$$(lsof -t -i :8787 2>/dev/null); \
	if [ -z "$$pids" ]; then \
		echo "No process is using port 8787"; \
	else \
		echo "Killing PID(s): $$pids on port 8787"; \
		kill $$pids; \
	fi
