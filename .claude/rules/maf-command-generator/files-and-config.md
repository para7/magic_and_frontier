---
paths:
  - "maf-command-generator/app/files/**/*.go"
  - "maf-command-generator/config/**"
  - "maf-command-generator/savedata/**"
---

# files 層・設定・データファイルの規約

## files パッケージ

- `JsonStore[T]`: `{ "entries": [...] }` 形式の JSON ファイルを読み書きする汎用ストア
- `MafConfig`: 全データファイルパスのハードコード設定（`LoadConfig()` で生成）
- `ExportSettings` / `ExportPaths`: `config/export_settings.json` から読み込むエクスポート先パス設定

## savedata ディレクトリ

各エンティティの JSON データ格納。ファイル名は `MafConfig` で定義:

- `grimoire.json`, `item.json`, `passive.json`, `bow.json`, `enemy_skill.json`
- `enemy.json`, `spawn_table.json`, `treasure.json`, `loottables.json`

## 注意

- `savedata/` や `config/` のデータを直接編集することは想定していない
- データの変更は MafEntity の CRUD メソッド経由で行う
