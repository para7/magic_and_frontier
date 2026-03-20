# Domain Migration Runbook

## Objective
- Recreate `@maf/domain` capabilities in Go with compatible data structures and core validations for happy paths.

## Input References
- `../tools/domain/src/**`
- `../tools/frontend/src/app/types.ts`
- `../tools/server/src/app.ts` (validation behavior reference)

## Output Location
- Recommended package roots:
1. `app/internal/domain/common`
2. `app/internal/domain/items`
3. `app/internal/domain/grimoire`
4. `app/internal/domain/skills`
5. `app/internal/domain/enemyskills`
6. `app/internal/domain/treasures`
7. `app/internal/domain/enemies`

## Task List
1. Create shared types/helpers
- `FieldErrors` map type.
- Common normalize functions (`trim`, newline normalization).
- UUID validation helper.
- Numeric parse/range helpers.

2. Define canonical entities
- `ItemEntry`
- `GrimoireEntry` and variant type
- `SkillEntry`
- `EnemySkillEntry` and trigger enum
- `DropRef`
- `SpawnRule`
- `EnemyEntry`
- `TreasureEntry`

3. Implement per-entity validation + normalization (happy path first)
- Required fields and min/max lengths.
- Required numeric fields and ranges.
- Optional numeric fields.
- Reference validation hooks via injected id-sets.

4. Implement upsert/delete helpers
- Generic upsert by `id` with mode result (`created`/`updated`).
- Generic delete by `id` with found/not-found result.

5. Implement state abstractions
- `EntryState[T]` as `{ entries: [] }`.
- Item and grimoire state wrappers per compatibility contract.

6. Implement domain result contracts
- Save result union equivalent:
1. success: `{ ok: true, entry, mode }`
2. failure: `{ ok: false, fieldErrors, formError }`
- Delete result equivalent:
1. success: `{ ok: true, deletedId }`
2. failure: `{ ok: false, formError, code? }`

7. Implement tests (happy path baseline)
- Save valid entry for each entity.
- Upsert existing id returns `updated`.
- Delete existing id returns `ok`.
- State load/save roundtrip for JSON shape compatibility.

## API/Type Decisions Locked
- `updatedAt` is ISO8601 string in UTC.
- `DropRef.kind` allowed values: `item`, `grimoire`.
- `EnemySkill.trigger` allowed values: `on_spawn`, `on_hit`, `on_low_hp`, `on_timer`.

## Completion Checklist
- [ ] All six entities compile in Go.
- [ ] Shared helper package reused by all validators.
- [ ] Domain tests pass (`go test ./...`).
- [ ] JSON structure remains compatible with existing savedata.
