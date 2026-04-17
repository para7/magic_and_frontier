---
paths:
  - "maf-command-generator/app/files/**/*.go"
  - "maf-command-generator/config/**"
  - "maf-command-generator/savedata/**"
---

# files 層・設定・データファイルの規約

## files パッケージ

- `JsonStore[T]`: 指定ディレクトリ配下の全 `*.json`（`{ "entries": [...] }` 形式）を走査してマージロードする汎用ストア
- `MafConfig`: 全エンティティの savedata パス・エクスポート設定パス・バニラ loot table ルートをハードコードで持つ設定（`LoadConfig()` で生成）。主なフィールド: `ItemStatePath`, `GrimoireStatePath`, `PassiveStatePath`, `BowStatePath`, `EnemySkillStatePath`, `EnemyStatePath`, `SpawnTableStatePath`, `TreasureStatePath`, `LootTableSourceRoot`, `ExportSettingsPath`, `MinecraftLootTableRoot`
- `ExportSettings` / `ExportPaths`: `config/export_settings.json` から読み込むエクスポート先パス設定

## savedata ディレクトリ

各エンティティは `savedata/{name}/` ディレクトリに格納され、配下の全 `*.json` ファイル（規約として `entity.json`）の `entries` がロード時にマージされる。`MafConfig` で定義されるパス:

- `savedata/grimoire/`, `savedata/item/`, `savedata/passive/`, `savedata/bow/`
- `savedata/enemy_skill/`, `savedata/enemy/`, `savedata/spawn_table/`, `savedata/treasure/`
- `savedata/loot_table/{namespace}/...` — Treasure エクスポートの入力（名前空間別 loot table JSON）

## 注意

- `savedata/` や `config/` のデータを直接編集することは想定していない
