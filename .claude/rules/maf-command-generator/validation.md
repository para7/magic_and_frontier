---
paths:
  - "maf-command-generator/app/domain/custom_validator/**/*.go"
  - "maf-command-generator/app/domain/model/validation.go"
  - "maf-command-generator/app/domain/model/validate_helpers.go"
---

# バリデーションシステム

## ValidationError（`model/validation.go`）

統一エラー型: Entity（エンティティ名）, ID, Field（JSON タグ名）, Tag, Param。

## custom_validator パッケージ

`go-playground/validator/v10` のラッパー。JSON タグ名をフィールド名として使用する設定済み。

### カスタムタグ

| タグ | 内容 |
|------|------|
| `trimmed_required` | `NormalizeText()` 後に空でないこと |
| `trimmed_min=N` | トリム後のルーン数が N 以上 |
| `trimmed_max=N` | トリム後のルーン数が N 以下 |
| `trimmed_oneof=a b c` | トリム後の値がいずれかに一致 |
| `maf_slug_id` | MAF 用スラッグ ID 形式の検証 |

### NormalizeText

CRLF/CR → LF 統一 + 前後空白除去。全入力テキストに適用。

### 再エクスポート

`ValidationErrors`, `FieldError`, `NewValidationError()`, `FormatValidationError()` を提供し、各パッケージが `go-playground/validator` を直接 import する必要をなくしている。

## validate_helpers.go

- `ValidateDropRefs()`: DropRef スライスの参照先存在確認（item/grimoire/passive/minecraft_item）+ CountMin <= CountMax チェック。passive の場合は Slot 必須チェックも行う
- `ValidateMafLootPools()`: `maf:item` / `maf:grimoire` / `maf:passive` を含む loot pool（任意の JSON マップ）を走査して参照先・slot・count を検証する
- `ValidateEquipmentSlots()`: Equipment 6 スロットの参照先存在確認
- `IsNamespacedResourceID()`, `IsSafeNamespacedResourcePath()`, `NormalizeResourcePath()`: リソースID/パスの検証・正規化
- `ParseLootEntryCount()`, `ParseLootEntrySlot()`: loot entry JSON から count/slot を取り出すヘルパー（バリデーション + convert で共有）

> ID 重複チェックは共通ヘルパーではなく、各エンティティの `ValidateAll` 内で `seenIDs map[string]bool` を使って自前で実装する。

## 共通バリデーションとドメイン固有バリデーションの境界

- `validate_helpers.go` には「全エンティティで再利用できる機械的ルール」だけを置く
- 特定ドメインの業務知識（例: passive の `generate_grimoire` が true でないと loot 参照不可）は各エンティティの `ValidateRelation()` で判定する
- export/convert 側で業務ルールエラーを返す設計を避け、`make run/validate` の段階で検出できるようにする
