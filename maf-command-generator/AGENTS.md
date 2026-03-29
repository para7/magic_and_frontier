# project info

`maf-command-editor` でデータ入力し、`generated` に出力する設計。`generated` 配下に問題があっても直接編集せず、生成元ロジックを修正する。

開発初期段階のため、過去互換より設計改善を優先する。

# Layering

現在の依存方向は固定で、逆流させない。

`app/main.go` -> `app/cli/*` -> `app/domain/master` -> `app/domain/model/*` / `app/domain/export/*` -> `app/files/*`

- `main`: サブコマンド分岐のみ。
- `cli`: 入出力と終了コード管理のみ。ドメインロジックは持たない。
- `master`: モデル間連携のハブ。model/export 用インターフェース実装を集約。
- `model`: エンティティ単位のCRUD/検証/永続化。
- `export`: 変換と出力。`DBMaster` 経由で読むだけ。
- `files`: JSONストアと設定ロードなどのI/Oユーティリティ。

# Interface map

## model層

`app/domain/model/interfaces.go`

- `type DBMaster interface { Has* }`
  - 各エンティティのリレーション検証で使う存在確認API。
- `type MafEntity[T any] interface`
  - `ValidateJSON(data, mas)`
  - `Create/Update/Delete(data/id, mas)`
  - `Save/Load()`
  - `ValidateAll(mas)`
  - `Find(id), GetAll()`

## export層

`app/domain/export/interfaces.go`

- `type DBMaster interface`
  - `GetGrimoireByID(id)`
  - `ListGrimoires()`

export層は model層の `MafEntity` に直接依存しない。必要最小限の読取インターフェースだけを追加する。

## master層

`app/domain/master/master.go`

- `DBMaster` が model/export 両方の `DBMaster interface` を実装する。
- `NewDBMaster(cfg)` で全エンティティを初期化し `Load()` 済みの状態を返す。
- `List*` は内部スライスを直接返さず defensive copy を返す。

# MafEntity 実装方針

各エンティティ実装ファイル（例: `app/domain/model/grimoire/entity.go`）は、以下の責務順で関数を並べる。

1. `NewXxxEntity(path)` コンストラクタ
2. `ValidateJSON`（`ValidateStruct` + `ValidateRelation` 合成）
3. `ValidateStruct`
4. `ValidateRelation`
5. `Create`
6. `Update`
7. `Delete`
8. `Save`
9. `Load`
10. `ValidateAll`
11. `Find`
12. `GetAll`

実装ルール:

- 構造体は `store files.JsonStore[T]` と `data []T` を持つ。
- バリデーションエラーは `[]model.ValidationError` で返す。
- `Create/Update` は先に `ValidateJSON` を実行し、エラー時は先頭エラーを `error` に変換して返す。
- `ValidateRelation` で他エンティティ存在確認が必要なら `mas.Has*` を使う。
- `Find` は `(T, bool)`、`GetAll` は現状データを返す。

# export 実装方針

- 変換専用の純粋関数は `convert.go` に置く（例: `grimoireToBook`）。
- 生成オブジェクト構築は `Build*Artifacts` に置く（副作用なし）。
- ファイル書き込みは `Write*Artifacts` に置く（副作用あり）。
- パス解決と設定読込は `ExportDatapack` で組み立てる。

# CLI 実装方針

- `Editor/Validate/Export` は `int` を返し、`main` が `os.Exit` する。
- `Export` は必ず `Validate` を先に実行し、失敗時は出力処理を中断する。
- 表示フォーマット調整は `cli` 側で行い、ドメイン層へ持ち込まない。

# 新規エンティティ追加手順

1. `app/domain/model/<entity>/types.go` に型定義（json/validateタグ込み）。
2. `app/domain/model/<entity>/entity.go` に `MafEntity` 実装。
3. `app/domain/model/interfaces.go` の `DBMaster` に必要な `Has*` を追加。
4. `app/domain/master/master.go` にフィールド追加、`NewDBMaster` で `Load()`、`Has*` 実装、必要なら `Get/List` 追加。
5. 必要に応じて `app/domain/export/interfaces.go` と export 実装を拡張。
6. `savedata/*.json` とテストを追加し、`make check` を通す。

# Repository Guidelines

- build check command: `make check`
