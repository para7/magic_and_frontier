---
paths:
  - "maf-command-generator/app/files/**/*.go"
  - "maf-command-generator/config/**"
  - "maf-command-generator/savedata/**"
---

# files 層・設定・データファイルの規約

## files パッケージ

- `JsonStore[T]`: 指定ディレクトリ配下の全 `*.json`（`{ "entries": [...] }` 形式）を走査してマージロードする汎用ストア
- `MafConfig`: 全エンティティの savedata パス・エクスポート設定パス・バニラ loot table ルートをハードコードで持つ設定（`LoadConfig()` で生成）。主なフィールド: `ItemStatePath`, `ActiveStatePath`, `PassiveStatePath`, `BowStatePath`, `EnemySkillStatePath`, `EnemyStatePath`, `SpawnTableStatePath`, `LootTableSourceRoot`, `ExportSettingsPath`, `MinecraftLootTableRoot`
- `ExportSettings` / `ExportPaths`: `config/export_settings.json` から読み込むエクスポート先パス設定。キー一覧: `activeEffect`, `activeGive`, `itemGive`, `passiveEffect`, `passiveApply`, `passiveGive`, `bowFlying`, `bowGround`, `enemy`, `enemySkill`, `enemyLoot`, `spawnTable`（デフォルト `generated/enemy/replace`）

## savedata ディレクトリ

各エンティティは `savedata/{name}/` ディレクトリに格納され、配下の全 `*.json` ファイル（規約として `entity.json`）の `entries` がロード時にマージされる。`MafConfig` で定義されるパス:

- `savedata/active/`, `savedata/item/`, `savedata/passive/`, `savedata/bow/`
- `savedata/enemy_skill/`, `savedata/enemy/`, `savedata/spawn_table/`
- `savedata/loot_table/{namespace}/...` — Treasure エクスポートの入力（名前空間別 loot table JSON）

### spawn_table の JSON 形式（例外）

`savedata/spawn_table/*.json` は `entries` 配列を持つが構造が独特:

```json
{
  "coordinates": { "dimension", "minDistance", "maxDistance", "minX/maxX", "minY/maxY", "minZ/maxZ" },
  "entries": [ { "sourceMobType", "baseMob": { "weight", "attributes": {"hp","attack","defense","moveSpeed"} }, "replacements": [...] } ]
}
```

- `coordinates` はファイル単位。ロード時に各 entry へ複製される
- `id` はオプション。省略時はファイル名 + sourceMobType + index から自動生成
- 他エンティティとは違い `JsonStore[T]` ではなく `spawntable.Entity.Load()` が独自パース

## 注意

- `savedata/` や `config/` のデータを直接編集することは想定していない
