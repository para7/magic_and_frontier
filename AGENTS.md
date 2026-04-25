# Codex Working Guide

詳細ルールは `./.claude/rules/` を正本とする。編集対象パスに応じて該当ドキュメントを読むこと。

## 参照ルーティング

| 対象パス (glob) | 参照ルール |
|----------------|-----------|
| `datapacks/magic_and_frontier/**` | `./.claude/rules/datapacks/always.md` |
| `maf-command-generator/app/**/*.go` | `./.claude/rules/maf-command-generator/architecture.md` |
| `maf-command-generator/app/{main.go,cli/**}` | `./.claude/rules/maf-command-generator/cli.md` |
| `maf-command-generator/app/domain/export/**` | `./.claude/rules/maf-command-generator/export.md` |
| `maf-command-generator/app/domain/master/**` | `./.claude/rules/maf-command-generator/master.md` |
| `maf-command-generator/app/domain/model/**/entity{,_test}.go` | `./.claude/rules/maf-command-generator/model-entity.md` |
| `maf-command-generator/app/domain/model/**/types.go` | `./.claude/rules/maf-command-generator/model-types.md` |
| `maf-command-generator/app/domain/model/item/components{,_test}.go` | `./.claude/rules/maf-command-generator/item-components.md` |
| `maf-command-generator/app/domain/{custom_validator,model/validation.go,model/validate_helpers.go}` | `./.claude/rules/maf-command-generator/validation.md` |
| `maf-command-generator/{app/files,config,savedata}/**` | `./.claude/rules/maf-command-generator/files-and-config.md` |
| `maf-command-generator/{app/minecraft,minecraft}/**` | `./.claude/rules/maf-command-generator/minecraft.md` |

ルールパスは `./.claude/rules/` からの相対。複数パターンに一致する場合はすべて読み、具体的なパスを優先する。

## 運用ルール

- ルール追加・変更は `./.claude/rules/` 側を更新する（サブディレクトリへ `AGENTS.md` を増やさない）
