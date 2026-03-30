---
paths:
  - "app/domain/export/**/*.go"
---

# export 層の規約

export 層は `export.DBMaster` 経由でデータを読み取り、Minecraft データパック（`.mcfunction` / loot table JSON）を生成する。`model.MafEntity` に直接依存しない。

## ファイル責務

| ファイル | 責務 |
|---------|------|
| `interfaces.go` | `export.DBMaster` インターフェース定義（読取専用） |
| `export.go` | `ExportDatapack`: オーケストレーション + パス解決 |
| `convert.go` | 純粋変換関数（`grimoireToBook`, `buildDropLootPool`, `toEnemyFunctionLines` 等） |
| `io.go` | `writeFunctionFile`, `writeJSON`: ファイル書き込みユーティリティ |
| `grimoire_effect.go` | `BuildGrimoireArtifacts` / `WriteGrimoireArtifacts` / `WriteGrimoireDebugArtifacts` |
| `enemyskill.go` | `BuildEnemySkillArtifacts` / `WriteEnemySkillArtifacts` |
| `enemy.go` | `BuildEnemyArtifacts` / `WriteEnemyArtifacts` |

## 設計原則

- 変換専用の純粋関数は `convert.go` に置く（例: `grimoireToBook`）
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
