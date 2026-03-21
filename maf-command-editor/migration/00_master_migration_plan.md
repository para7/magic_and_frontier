# Go Migration Master Plan (`../tools` -> `tools2`)

## Goal
- Migrate `domain`, `server`, and `frontend` from JavaScript/TypeScript in `../tools` to Go in `tools2`.
- Keep `savedata/*.json` and `server/config/export-settings.json` compatibility.
- Keep REST endpoints under `/api/*` during migration.
- Rebuild UI as server-rendered pages using `templ + htmx`.

## Scope
- Features included:
1. `items`
2. `grimoire`
3. `skills`
4. `enemy-skills`
5. `treasures`
6. `enemies`
7. `POST /api/save`

## Non-goals (for this phase)
- Full parity with all existing JS validation edge cases.
- Exhaustive failure-path test migration.
- Keeping Angular frontend running in parallel.

## Migration Order
- Package order:
1. `domain`
2. `server`
3. `frontend`
- Feature order:
1. `items`
2. `grimoire`
3. `skills`
4. `enemy-skills`
5. `treasures`
6. `enemies`
7. `save`

## Compatibility Contracts
- API routes must remain:
1. `GET/POST/DELETE /api/items`
2. `GET/POST/DELETE /api/grimoire`
3. `GET/POST/DELETE /api/skills`
4. `GET/POST/DELETE /api/enemy-skills`
5. `GET/POST/DELETE /api/enemies`
6. `GET/POST/DELETE /api/treasures`
7. `POST /api/save`
- Storage files must remain readable/writable with current key structures.
- `tools2` existing form app may be used as implementation reference, then replaced by migrated app pages.

## Architecture Target (Go)
- `app/internal/domain`: entity models, normalization, validation, usecases.
- `app/internal/store`: JSON repositories and config loading.
- `app/internal/httpapi`: `/api/*` handlers.
- `app/internal/web`: SSR page handlers (`templ + htmx` partial updates).
- `app/views`: layout + per-feature pages/components.

## Milestones
1. Domain foundation
- Shared parse/validate helpers.
- Generic entry repository (`{ entries: [] }`) and item/grimoire specialized states.
- Time provider abstraction for `updatedAt`.

2. API foundation
- Health endpoint.
- CRUD endpoints for `items` and `grimoire`.
- Unified error/response shape.

3. Domain/API completion
- Add `skills`, `enemy-skills`, `treasures`, `enemies`.
- Reference integrity checks (`itemId`, `enemySkillIds`, loot refs).

4. Save/export
- Add `POST /api/save`.
- Read export settings; generate output payload and files.

5. SSR frontend completion
- List/create/update/delete pages for all six features.
- Save button flow and success/failure notification.

6. Cleanup
- Remove obsolete demo-form routes and templates from `tools2`.
- Update Makefile targets and test command defaults if needed.

## Acceptance Criteria
1. `make test` passes with happy-path tests for each feature.
2. `/api/*` endpoints respond with expected JSON shape on success.
3. Existing savedata files can be loaded without manual migration.
4. SSR UI supports create/list/delete flow for all six features.
5. `POST /api/save` succeeds with valid settings and data.

## Execution Rules for Codex
- Implement one feature at a time following feature order.
- For each feature:
1. Add/adjust domain types and validation.
2. Add repository/usecase operations.
3. Add API handlers and tests.
4. Add SSR views/handlers and smoke tests.
- Do not break previously completed feature tests.
- Keep each commit small and feature-scoped.
