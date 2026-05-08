---
paths:
  - "maf-command-generator/app/domain/export/**/*.go"
---

# export 層の規約

export 層は `export.DBMaster` 経由でデータを読み取り、Minecraft データパック（`.mcfunction` / loot table JSON）を生成する。`model.MafEntity` に直接依存しない。

## ファイル責務

| ファイル | 責務 |
|---------|------|
| `interfaces.go` | `export.DBMaster` インターフェース定義（読取専用） |
| `export.go` | `ExportDatapack`: オーケストレーション + パス解決 |
| `convert/` | 純粋変換関数群（`export_convert` パッケージ） |
| `io.go` | `writeFunctionFile`, `writeJSON`: ファイル書き込みユーティリティ |
| `master_lookup.go` | `buildMasterEntityLookups`: master から items/actives/passives/bows をまとめてロードし ID マップ化（内部） |
| `active_effect.go` | `BuildActiveArtifacts` / `WriteActiveArtifacts` / `WriteActiveDebugArtifacts` |
| `item.go` | `BuildItemArtifacts` / `WriteItemArtifacts` |
| `passive.go` | `BuildPassiveArtifacts` / `WritePassiveArtifacts` |
| `bow.go` | `BuildBowArtifacts` / `WriteBowArtifacts` |
| `enemyskill.go` | `BuildEnemySkillArtifacts` / `WriteEnemySkillArtifacts` |
| `enemy.go` | `BuildEnemyArtifacts` / `WriteEnemyArtifacts` |
| `spawntable.go` | `BuildSpawnTableArtifacts` / `WriteSpawnTableArtifacts`（バニラモブを重み付き確率でカスタムエネミーに置換。`main.mcfunction` がディスパッチャ、各 `{id}.mcfunction` が `maf_vh_rand` スコアで分岐。`BaseMob.Attributes` があれば残存バニラモブに `data merge entity` で Health / Attributes NBT を書き込む） |
| `treasure.go` | `BuildTreasureArtifacts` / `WriteTreasureArtifacts`（`LootTableSourceRoot` 配下の JSON を走査し、`maf:*` エントリ解決 + `minecraft:` 名前空間はバニラとマージ） |

## convert/ サブパッケージ（export_convert）

エンティティデータを `.mcfunction` コマンド文字列や loot table JSON に変換する純粋関数群。主な変換:

- アクティブ/パッシブ → 本アイテム SNBT（`spellBookModel` 共通基盤、パッシブは詠唱 200tick / MP 10 固定）
- アイテム → give コマンド SNBT（custom_data・コンポーネント・エンチャント組み立て）
- `maf:*` loot エントリ → バニラ互換 loot entry 解決・マージ
- エネミー → summon NBT + function 行

`export` パッケージを逆インポートしてはならない（循環防止）。import alias: `ec "maf_command_editor/app/domain/export/convert"`

## 設計原則

- 変換専用の純粋関数は `convert/` サブパッケージに置く
- 生成オブジェクト構築は `Build*Artifacts` に置く（副作用なし）
- ファイル書き込みは `Write*Artifacts` に置く（副作用あり）
- パス解決と設定読込は `ExportDatapack` で組み立てる
- `convert/` は validate の代替をしない。ドメイン制約違反は model の `ValidateRelation` で検出する
- 生成条件の変更で出力対象が減る場合は、`Write*Artifacts` 側で stale ファイルを削除して前回生成物を残さない

## export.DBMaster インターフェース

新しいエンティティのエクスポートが必要になったら、`export.DBMaster` に必要最小限の読取メソッドだけを追加する。

## 出力構成

- `.mcfunction` → `{outputRoot}/data/maf/function/{logicalDir}/`
- loot table JSON → `{outputRoot}/data/maf/loot_table/{logicalDir}/`
- パスは `config/export_settings.json` で設定
- `passive/bow/` の出力先のみ `export.go` 内でハードコード（`generated/passive/bow`）。`export_settings.json` では未定義

## Enemy ドロップモード

- `replace`: カスタムプールのみで loot table を構築
- `append`: バニラ loot table（`minecraft/1.21.11/loot_table/` 配下）を読み込み、カスタムプールをマージ

## Treasure エクスポート

`BuildTreasureArtifacts` は `savedata/loot_table/{namespace}/...` を走査し、`maf:item` / `maf:active` / `maf:passive` エントリを vanilla 互換の loot entry に解決する。`minecraft` 名前空間のファイルは `minecraft/1.21.11/loot_table/` 配下のバニラ loot table を読み込み、カスタムプールを追記する（`append` 相当）。出力先は `{outputRoot}/data/{namespace}/loot_table/{relPath}.json`。
