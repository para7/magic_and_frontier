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
| `grimoire_effect.go` | `BuildGrimoireArtifacts` / `WriteGrimoireArtifacts` / `WriteGrimoireDebugArtifacts` |
| `item.go` | `BuildItemArtifacts` / `WriteItemArtifacts` |
| `passive.go` | `BuildPassiveArtifacts` / `WritePassiveArtifacts` |
| `bow.go` | `BuildBowArtifacts` / `WriteBowArtifacts` |
| `enemyskill.go` | `BuildEnemySkillArtifacts` / `WriteEnemySkillArtifacts` |
| `enemy.go` | `BuildEnemyArtifacts` / `WriteEnemyArtifacts` |

## convert/ サブパッケージ（export_convert）

| ファイル | 公開関数 | 内容 |
|---------|---------|------|
| `json.go` | `JsonString`, `ToCountValue`, `ToWeight` | JSON文字列化・float変換・loot カウント/重みヘルパー |
| `grimoire.go` | `GrimoireToBook` | グリモア → 本アイテム SNBT 変換 |
| `passive.go` | `PassiveToBook` | パッシブ → 本アイテム SNBT 変換 |
| `book.go` | `spellBookModel` | グリモア/パッシブ共通の本アイテム生成モデル |
| `item.go` | なし | item give コマンド・コンポーネント・エンチャント（内部） |
| `loottable.go` | `BuildDropLootPool`, `MergeLootTablePools` | loot pool 構築・マージ |
| `enemy.go` | `ToEnemyFunctionLines` | エネミー summon NBT と function 行生成 |

このサブパッケージは `export` パッケージを逆インポートしてはならない（循環防止）。
import alias は `ec "maf_command_editor/app/domain/export/convert"` を使用。

## 設計原則

- 変換専用の純粋関数は `convert/` サブパッケージに置く
- 生成オブジェクト構築は `Build*Artifacts` に置く（副作用なし）
- ファイル書き込みは `Write*Artifacts` に置く（副作用あり）
- パス解決と設定読込は `ExportDatapack` で組み立てる

## export.DBMaster インターフェース

新しいエンティティのエクスポートが必要になったら、`export.DBMaster` に必要最小限の読取メソッドだけを追加する。

## 出力構成

- `.mcfunction` → `{outputRoot}/data/maf/function/{logicalDir}/`
- loot table JSON → `{outputRoot}/data/maf/loot_table/{logicalDir}/`
- パスは `config/export_settings.json` で設定

## Enemy ドロップモード

- `replace`: カスタムプールのみで loot table を構築
- `append`: バニラ loot table（`minecraft/1.21.11/loot_table/` 配下）を読み込み、カスタムプールをマージ
