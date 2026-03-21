# Server Migration Runbook

## Objective
- Rebuild `@maf/server` behavior in Go, preserving route contracts and storage compatibility.

## Input References
- `../tools/server/src/app.ts`
- `../tools/server/src/config.ts`
- `../tools/server/src/repositories/json-file-repositories.ts`
- `../tools/server/test/*.test.ts`

## Output Location
- Recommended package roots:
1. `app/internal/config`
2. `app/internal/store`
3. `app/internal/httpapi`
4. `app/internal/export`

## Task List
1. Implement config loader
- Env vars:
1. `PORT`
2. `ITEM_STATE_PATH`
3. `GRIMOIRE_STATE_PATH`
4. `SKILL_STATE_PATH`
5. `ENEMY_SKILL_STATE_PATH`
6. `ENEMY_STATE_PATH`
7. `TREASURE_STATE_PATH`
8. `EXPORT_SETTINGS_PATH`
- Candidate path fallback logic compatible with current JS behavior.

2. Implement JSON repositories
- Generic read/write JSON helpers.
- `ENOENT` -> default empty state behavior.
- Item/grimoire specialized repositories.
- Generic entry repository for `{ entries: [] }`.

3. Implement API handlers
- `GET /health`
- CRUD:
1. `/api/items`
2. `/api/grimoire`
3. `/api/skills`
4. `/api/enemy-skills`
5. `/api/enemies`
6. `/api/treasures`
- `POST /api/save`

4. Implement response mapping
- Success/failure JSON shape aligned to existing client expectations.
- HTTP status mapping:
1. Validation failure -> `400`
2. Not found on delete -> `404`
3. Success -> `200`

5. Implement cross-entity validation hooks
- `skills.itemId` must exist.
- `enemies.enemySkillIds[]` must exist.
- `treasures.lootPools[]` refs must exist.
- `enemies.dropTable[]` refs must exist when provided.
- Prevent deleting referenced `enemy-skill` with `REFERENCE_ERROR`.

6. Implement save/export bridge
- Read current states and export settings.
- Return save result payload with generated counts/message.
- Keep behavior stable for minimum happy path execution.

7. Integrate router into `app/main.go`
- Replace demo-form routes progressively.
- Keep middleware: recover + request logging.
- Move the old form demo into a non-build archive path and document that it is retained only as plan 3 reference material.

8. Add server tests (happy path baseline)
- Health endpoint returns `{ ok: true }`.
- CRUD success path for all six entities.
- Reference-resolved save for dependent entities.
- `POST /api/save` success with temp output dir fixture.

## Completion Checklist
- [ ] All `/api/*` routes implemented in Go.
- [ ] Existing savedata files load without manual conversion.
- [ ] Server tests pass for happy paths.
- [ ] Demo-form-only endpoints removed or isolated.
