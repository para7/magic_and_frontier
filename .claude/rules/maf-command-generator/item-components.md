---
paths:
  - "maf-command-generator/app/domain/model/item/components.go"
  - "maf-command-generator/app/domain/model/item/components_test.go"
---

# Item Components

`BuildItemComponents` は Item データから Minecraft 1.20.5+ の item component SNBT 文字列を生成する。

## 入力

Item の `Minecraft.Components`（`map[string]string`）から各コンポーネントを取得。
キーは `minecraft:*` 形式の namespace 付きコンポーネント名、値は SNBT 文字列。

## 変換処理

1. `NormalizeComponents` でキー/値のトリム・バリデーション（空キー・namespace なしキー・空値をエラー）
2. キーをアルファベット順にソート
3. 各コンポーネントを `"key":value` 形式で結合

## 出力形式

`{id:"minecraft:...",count:1,components:{...}}` 形式の SNBT 文字列。
