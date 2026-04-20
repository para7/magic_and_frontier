---
paths:
  - "maf-command-generator/app/domain/model/**/entity.go"
  - "maf-command-generator/app/domain/model/**/entity_test.go"
---

# MafEntity 実装規約

各エンティティは `model.MafEntity[T]` を実装する。

## 関数の配置順序

1. `NewXxxEntity(path)` コンストラクタ
2. `ValidateJSON`（`ValidateStruct` + `ValidateRelation` 合成）
3. `ValidateStruct`
4. `ValidateRelation`
5. `Load`
6. `ValidateAll`
7. `Find`
8. `GetAll`

## 実装ルール

- 構造体は `store files.JsonStore[T]` と `data []T` を持つ
- バリデーションエラーは `[]model.ValidationError` で返す
- `ValidateRelation` で他エンティティ存在確認が必要なら `mas.Has*` を使う
- `ValidateRelation` で固有業務ルールの判定に参照先詳細が必要な場合は、`model.DBMaster` の正式メソッド（例: `GetPassive` で `PassiveSnapshot` を取得）を使う
- `entity.go` 内で一時的なローカル interface を定義して責務を分岐しない。必要なら `model.DBMaster` を拡張して `master.DBMaster` に実装を追加する
- ID 重複チェックは `ValidateAll` 内で `map[string]bool` の `seenIDs` を使った走査で検出し、重複分を `Tag:"unique"` の `ValidationError` として追加する（汎用ヘルパーは未提供）
- `Find` は `(T, bool)`、`GetAll` は現在データをそのまま返す

### 例外: SpawnTable

`model/spawntable/entity.go` は `JsonStore` の汎用 `entries` 配列ではなく、独自の `{coordinates, entries}` 形式を読む。`Load()` は `savedata/spawn_table/*.json` を glob で走査し、ファイル単位の `coordinates`（`dimension`, `minDistance/maxDistance`, `minX~maxZ`）を各 `entries[i]` に複製してから `SpawnTable` に展開する。`id` が空のエントリはファイル名 + `sourceMobType` + index から `buildSpawnTableID` で自動付与される（小文字化・namespace 除去・非英数の `_` 置換）。新規エンティティでこの形式を真似る場合は `store` を使わず `*.json` 走査 + 独自パーサで実装する。

`ValidateRelation` では `BaseMob` 配下の属性範囲チェック（`validateBaseMobAttributes`: hp 1..100000, attack/defense/moveSpeed 0..100000）と `minDistance <= maxDistance` もここで実施する。

## バリデーション

- 構造体バリデーション: `custom_validator.Validate.Struct()` を使用
- カスタムタグ: `trimmed_required`, `trimmed_min`, `trimmed_max`, `trimmed_oneof`, `maf_slug_id`
- リレーションバリデーション: `model.ValidateDropRefs()`, `model.ValidateMafLootPools()`, `model.ValidateEquipmentSlots()` を活用
