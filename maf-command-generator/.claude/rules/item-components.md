---
paths:
  - "app/domain/model/item/components.go"
  - "app/domain/model/item/components_test.go"
---

# Item Components

`BuildItemComponents` は Item データから Minecraft 1.20.5+ の item component SNBT 文字列を生成する。

## 変換対象

エディタの各フィールドを対応する `minecraft:*` コンポーネントに変換:

- `CustomName` → `minecraft:custom_name`
- `Lore` → `minecraft:lore`
- `Enchantments` → `minecraft:enchantments`
- `Unbreakable` → `minecraft:unbreakable`
- `CustomModelData` → `minecraft:custom_model_data`
- `RepairCost` → `minecraft:repair_cost`
- `HideFlags` → `minecraft:tooltip_display`（ビットマスクから変換）
- `Potion*` → `minecraft:potion_contents`
- `AttributeModifiers` → `minecraft:attribute_modifiers`
- `CustomNBT` → 既出キーを除いた残りをそのまま追加

## 出力形式

`{id:"minecraft:...",count:1,components:{...}}` 形式の SNBT 文字列。
