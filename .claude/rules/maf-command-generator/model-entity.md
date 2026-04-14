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
5. `Create`
6. `Update`
7. `Delete`
8. `Save`
9. `Load`
10. `ValidateAll`
11. `Find`
12. `GetAll`

## 実装ルール

- 構造体は `store files.JsonStore[T]` と `data []T` を持つ
- バリデーションエラーは `[]model.ValidationError` で返す
- `Create`/`Update` は先に `ValidateJSON` を実行し、エラー時は先頭エラーを `error` に変換して返す
- `ValidateRelation` で他エンティティ存在確認が必要なら `mas.Has*` を使う
- `ValidateRelation` で固有業務ルールの判定に参照先詳細が必要な場合は、`model.DBMaster` の正式メソッド（例: `GetPassive`）を使う
- `entity.go` 内で一時的なローカル interface を定義して責務を分岐しない。必要なら `model.DBMaster` を拡張して `master.DBMaster` に実装を追加する
- `ValidateRelation` で ID 重複チェックには `model.HasDuplicateID` を使う
- `ValidateAll` では `model.FilterUniqueIDErrorsPerEntityID` で ID 重複エラーを間引く
- `Find` は `(T, bool)`、`GetAll` は現在データをそのまま返す

## バリデーション

- 構造体バリデーション: `custom_validator.Validate.Struct()` を使用
- カスタムタグ: `trimmed_required`, `trimmed_min`, `trimmed_max`, `trimmed_oneof`, `maf_slug_id`
- リレーションバリデーション: `model.ValidateDropRefs()`, `model.ValidateEquipmentSlots()` を活用
