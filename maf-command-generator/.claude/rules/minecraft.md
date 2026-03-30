---
paths:
  - "app/minecraft/**/*.go"
  - "minecraft/**"
---

# minecraft パッケージの規約

バニラ Minecraft の loot table データを扱うユーティリティ。

## 主要機能

- `ListLootTables(root)`: 指定ディレクトリ以下の全 loot table を列挙
- `LoadLootTable(root, tablePath)`: `tablePath`（`minecraft:entities/zombie` 等）に対応する JSON を読み込み
- `Exists(root, tablePath)`: loot table ファイルの存在確認
- `FilePathForTable(root, tablePath)`: tablePath → 実ファイルパスの変換

## tablePath 形式

`minecraft:<category>/<name>` の namespace 形式（例: `minecraft:entities/zombie`）。
`minecraft/1.21.11/loot_table/` 配下にバニラデータのスナップショットを保持。

## 用途

- `model.DBMaster.HasMinecraftLootTable()`: Treasure の参照先存在確認
- `export.BuildEnemyArtifacts()`: append モードでバニラ loot table をベースにカスタムプールをマージ
