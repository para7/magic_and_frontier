---
name: add-maf-item
description: maf-command-generatorのsavedataにアイテムを追加するスキル。item, アイテム, 武器, 防具, 装備, 剣, 弓, ヘルメット, チェストプレート などの新規追加が出たら使う。savedata/item/claude.json への書き込みを担う。
disable-model-invocation: true
---

# アイテム追加

**対象ファイル:** `maf-command-generator/savedata/item/claude.json`

## 手順

1. `claude.json` が存在すれば読み込む。なければ `{"entries": []}` として扱う
2. `entries` に新エントリを追加
3. ファイルに書き戻す
4. 追加後: `cd maf-command-generator && make run/validate` で検証

## JSONスキーマ

```json
{
  "entries": [
    {
      "id": "item_id",            // 必須: 小文字・アンダースコア・ハイフンのみ、ユニーク
      "maf": {                    // 必須(空オブジェクト可): maf固有データ
        "grimoireId": "spell01",  // 任意: バインドするグリモアID
        "passiveId": "regen",     // 任意: バインドするパッシブID
        "bowId": "bow_skill01",   // 任意: バインドする弓パッシブID
        "maxmp": 20               // 任意: MaxMP増減値(負数可)
      },
      "minecraft": {              // 必須: Minecraftアイテムデータ
        "itemId": "minecraft:diamond_sword",  // 必須: アイテムID
        "components": {           // 必須(空オブジェクト可): コンポーネント
          "minecraft:custom_name": "'{\"text\":\"表示名\"}'",
          "minecraft:lore": "['{\"text\":\"説明文\"}']",
          "minecraft:enchantments": "{\"minecraft:sharpness\":5}"  // 任意
        }
      }
    }
  ]
}
```

## mafフィールドの組み合わせ

| 用途 | 設定するフィールド |
|------|-----------------|
| 通常アイテム | `maf: {}` |
| 魔法武器 | `grimoireId` |
| パッシブ武器 | `passiveId` |
| 弓スキル武器 | `bowId` |
| MP装備(増) | `maxmp: 20` |
| MP装備(減) | `maxmp: -20` |
| 複合 | 複数フィールド同時設定可 |

## custom_nameとlore書式

```
"minecraft:custom_name": "'{\"text\":\"アイテム名\"}'"
"minecraft:lore": "['{\"text\":\"1行目\"}', '{\"text\":\"2行目\"}']"
```

## よく使うitemId

- 武器: `minecraft:diamond_sword`, `minecraft:iron_sword`
- 弓: `minecraft:bow`, `minecraft:crossbow`
- 防具: `minecraft:diamond_helmet`, `minecraft:diamond_chestplate`, `minecraft:diamond_leggings`, `minecraft:diamond_boots`
- その他: `minecraft:stick`, `minecraft:book`, `minecraft:stone`
