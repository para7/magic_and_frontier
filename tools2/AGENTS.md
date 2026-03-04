# Repository Guidelines

## Project Structure & Module Organization
This repository is a small Go web app (`module tools2`) with form validation and server-rendered templates.

- `main.go`: HTTP server entrypoint and route handlers (`GET /`, `POST /contact/submit`).
- `internal/form`: shared form state and error types.
- `internal/validation`: validation logic and unit tests.
- `views/*.templ`: `templ` component sources for page and form UI.
- `main_test.go`, `internal/**/*_test.go`: handler and package-level tests.
- `bin/`: local build output (`bin/tools2`).

## Build, Test, and Development Commands
Use `make` targets as the primary workflow:

- `make setup`: install `templ` and `staticcheck` tool dependencies.
- `make generate`: generate Go code from `views/*.templ`.
- `make format`: run `go fmt ./...`.
- `make lint`: run `go vet ./...` and `staticcheck ./...`.
- `make build`: compile binary to `bin/tools2`.
- `make test`: run all tests (`go test ./...`).
- `make run`: start local server on `:8080`.
- `make check`: run full CI-like pipeline (generate, tidy, format, lint, build, test).

## Coding Style & Naming Conventions
- Follow standard Go formatting; always run `make format`.
- Keep package names lowercase and concise (`form`, `validation`, `views`).
- Exported identifiers use `CamelCase`; unexported helpers use `camelCase`.
- Prefer table-driven tests for validation/business logic.
- Keep templates in `.templ` files and regenerate after edits (`make generate`).

## Testing Guidelines
- Framework: Go `testing` package with `httptest` for handler coverage.
- Test file naming: `*_test.go`; test function naming: `TestXxx`.
- Include both success and failure paths (example: `TestSubmitHandlerSuccess`, validation-error cases).
- Run `make test` locally before opening a PR.
- Always run `make check` before marking work as complete.

## Commit & Pull Request Guidelines
Recent history mixes concise commits (`save`) and descriptive commits (`chore: ...`). Prefer descriptive, imperative messages with optional scope:

- Example: `feat(validation): trim whitespace before checks`
- Example: `fix(handler): return 400 on invalid form parse`

For PRs, include:
- What changed and why.
- Commands run locally (for example: `make check`).
- Screenshots or HTML snippets when UI/template output changes.
