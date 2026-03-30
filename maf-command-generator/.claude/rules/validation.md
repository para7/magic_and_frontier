---
paths:
  - "app/domain/custom_validator/**/*.go"
  - "app/domain/model/validation.go"
  - "app/domain/model/validate_helpers.go"
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

### NormalizeText

CRLF/CR → LF 統一 + 前後空白除去。全入力テキストに適用。

### 再エクスポート

`ValidationErrors`, `FieldError`, `NewValidationError()`, `FormatValidationError()` を提供し、各パッケージが `go-playground/validator` を直接 import する必要をなくしている。

## validate_helpers.go

- `ValidateDropRefs()`: DropRef スライスの参照先存在確認（item/grimoire/minecraft_item）+ CountMin <= CountMax チェック
- `ValidateEquipmentSlots()`: Equipment 6 スロットの参照先存在確認
- `IsNamespacedResourceID()`, `IsSafeNamespacedResourcePath()`, `NormalizeResourcePath()`: リソースID/パスの検証・正規化
