# Magic and Frontier v2

Minecraft データパックの RPG 魔法システムと、それを生成する Go ツールのモノリポ。

## 構成

| ディレクトリ | 内容 |
|-------------|------|
| `datapacks/magic_and_frontier/` | メインデータパック（名前空間 `maf`） |
| `maf-command-generator/` | Go 製データパック生成ツール |
| `maf-command-generator/savedata/` | マスターデータ JSON |

## コマンド（maf-command-generator/）

- `make check` — フル検証（generate + tidy + format + lint + build + test）
- `make run/validate` — バリデーション実行
- `make run/export` — バリデーション + エクスポート実行

## ルール体系

パスベースの詳細ルールは `.claude/rules/` に配置。各ファイルの `paths:` frontmatter で自動ロードされる。

@AGENTS.md
