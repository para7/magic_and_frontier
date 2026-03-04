# Repository Guidelines

## Project Goal
For the near term, this repository’s primary goal is to migrate `../tools` to Go.

- Prioritize parity with behavior in `../tools` before adding new features.
- When changing logic, document migration intent and compatibility impact in PR notes.
- Prefer incremental, reviewable steps (for example: parser migration, then handler wiring, then template updates).

## Project Structure & Module Organization
This repository is a small Go web app with server-rendered templates.

- `app/`: main application code (`main.go`, handlers, internal domain logic).
- `app/internal/forms/`: form models and parsing/validation logic.
- `app/views/`: `templ` UI files, including `layout.templ` and form partials under `views/form/`.
- `migration/`: migration runbooks and planning docs.
- `bin/`: build output (`tools2` binary).
- Root files: `Makefile`, `go.mod`, `go.sum`, `mise.toml`.

Keep new domain logic in `app/internal/<feature>/` and keep HTTP wiring in `app/`.

## Build, Test, and Development Commands
Use `make` targets as the canonical workflow:

- `make setup`: install tool dependencies (`templ`, `staticcheck`).
- `make generate`: regenerate Go code from `.templ` files.
- `make format`: run `go fmt ./...`.
- `make lint`: run `go vet` and `staticcheck`.
- `make build`: build `bin/tools2`.
- `make test`: run all tests (`go test ./...`).
- `make run` / `make dev`: generate templates and run app on `http://localhost:8080`.
- `make check`: full CI-style pass (generate, tidy, format, lint, build, test).

## Coding Style & Naming Conventions
- Follow standard Go formatting and idioms; always run `make format`.
- Use `camelCase` for unexported names and `PascalCase` for exported identifiers.
- Keep package names short, lowercase, and purpose-driven (`forms`, not `form_utils`).
- Name template partials by purpose (for example, `fields_storage.templ`, `fields_message.templ`).
- Prefer small handler functions that delegate validation/model logic to `app/internal`.

## Testing Guidelines
- Use Go’s `testing` package; keep tests in `*_test.go` files next to implementation.
- Prefer table-driven tests for parser/validator logic in `app/internal/forms/`.
- Run `make test` locally before opening a PR.
- Run `make check` when touching templates, module dependencies, or request handling paths.

## Commit & Pull Request Guidelines
Recent history is mixed (`save`, Japanese summaries, and `chore:` style). For new work, use clear imperative subjects, optionally with a scope:

- `forms: validate storage field bounds`
- `templ: split server fields into partial`

PRs should include:

- concise problem/solution summary,
- linked issue or migration doc when applicable,
- test notes (`make test` / `make check` results),
- screenshots or rendered output notes for UI/template changes.
