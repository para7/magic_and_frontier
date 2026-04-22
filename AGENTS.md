# Codex Working Guide

このリポジトリでは、詳細ルールは `./.claude/rules/` を正本とし、Codex は編集対象パスに応じて該当ドキュメントを参照する。

## 参照ルーティング（glob相当）

- `datapacks/magic_and_frontier/**`
  - `./.claude/rules/datapacks/always.md`

- `maf-command-generator/app/**/*.go`
  - `./.claude/rules/maf-command-generator/architecture.md`

- `maf-command-generator/app/main.go`
- `maf-command-generator/app/cli/**/*.go`
  - `./.claude/rules/maf-command-generator/cli.md`

- `maf-command-generator/app/domain/export/**/*.go`
  - `./.claude/rules/maf-command-generator/export.md`

- `maf-command-generator/app/domain/master/**/*.go`
  - `./.claude/rules/maf-command-generator/master.md`

- `maf-command-generator/app/domain/model/**/entity.go`
- `maf-command-generator/app/domain/model/**/entity_test.go`
  - `./.claude/rules/maf-command-generator/model-entity.md`

- `maf-command-generator/app/domain/model/**/types.go`
  - `./.claude/rules/maf-command-generator/model-types.md`

- `maf-command-generator/app/domain/model/item/components.go`
- `maf-command-generator/app/domain/model/item/components_test.go`
  - `./.claude/rules/maf-command-generator/item-components.md`

- `maf-command-generator/app/domain/custom_validator/**/*.go`
- `maf-command-generator/app/domain/model/validation.go`
- `maf-command-generator/app/domain/model/validate_helpers.go`
  - `./.claude/rules/maf-command-generator/validation.md`

- `maf-command-generator/app/files/**/*.go`
- `maf-command-generator/config/**`
- `maf-command-generator/savedata/**`
  - `./.claude/rules/maf-command-generator/files-and-config.md`

- `maf-command-generator/app/minecraft/**/*.go`
- `maf-command-generator/minecraft/**`
  - `./.claude/rules/maf-command-generator/minecraft.md`

## 運用ルール

- 複数のパターンに一致する場合は、該当する `rules` をすべて読む。
- ルールの重複・競合がある場合は、より具体的なパスのルールを優先する。
- ルール追加や変更は、原則 `./.claude/rules/` 側を更新する。
- この方針により、サブディレクトリへ `AGENTS.md` を増やさない。
