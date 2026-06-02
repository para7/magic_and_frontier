# Plan: items.tsv から level1〜5.json 生成

## Context

`items.tsv` にアーマーの全レベルスペックが定義されている。これを `savedata/item/armors/levelN.json` に変換する。既存の `level1.json` は新データで置き換える。

## TSV列マッピング

| TSV列 | JSON フィールド | 備考 |
|-------|---------------|------|
| HP | `attribute_modifiers[type="max_health"].amount` | 0またはブランクはスキップ |
| MP | `maf.maxmp` | 0またはブランクはスキップ（maf全体省略） |
| 強度 | `attribute_modifiers[type="armor_toughness"].amount` | 0またはブランクはスキップ |
| Enchant | `minecraft:enchantments["minecraft:protection"]` | "-"または0またはブランクはスキップ |
| attribute | `attribute_modifiers[type="attack_damage"].amount` | "attack_damage+N"形式をパース、"-"はスキップ |

加えて、全アーマーピースにバニラ相当の `armor` 属性を常時追加（`attribute_modifiers` コンポーネントはデフォルト属性を置き換えるため必須）。

## バニラ armor 値（固定ベース）

| 素材 | helmet | chestplate | leggings | boots |
|------|--------|------------|----------|-------|
| leather | 1 | 3 | 2 | 1 |
| chainmail | 2 | 5 | 4 | 1 |
| iron | 2 | 6 | 5 | 2 |
| copper | 2 | 5 | 4 | 1 |（推定 - 実機確認要） |
| gold | 1 | 3 | 2 | 1 |
| diamond | 3 | 8 | 6 | 3 |
| netherite | 3 | 8 | 6 | 3 |

## 各レベルで生成するエントリ

素材ごとのレベル存在判定：

| 素材 | L1 | L2 | L3 | L4 | L5 |
|------|----|----|----|----|-----|
| netherite | ✗（none）| ✓ | ✓ | ✓ | ✓ |
| diamond | ✗（none）| ✓ | ✓ | ✓ | ✓ |
| iron | ✓ | ✓ | ✓ | ✓ | ✓ |
| chainmail | ✓ | ✓ | ✓ | ✓ | ✓ |
| copper | ✓ | ✓ | ✓ | ✓ | ✓ |
| gold | ✓ | ✓ | ✓ | ✓ | ✓ |
| leather | ✓ | ✓ | ✓ | ✓ | ✓ |

## TSVから読み取った各スペック

### Netherite（L2〜5）
| slot | L2 | L3 | L4 | L5 |
|------|----|----|----|----|
| helmet | HP=1,MP=1,強度=2,E=5 | HP=2,MP=1,強度=3,E=5 | HP=4,MP=2,強度=4,E=5 | HP=6,MP=5,強度=7,E=5 |
| chestplate | HP=2,MP=2,強度=3,E=5 | HP=2,MP=5,強度=5,E=8 | HP=5,MP=7,強度=7,E=9 | HP=12,MP=10,強度=16,E=10 |
| leggings | HP=2,MP=2,強度=3,E=5 | HP=2,MP=5,強度=5,E=8 | HP=5,MP=7,強度=7,E=9 | HP=12,MP=10,強度=16,E=10 |
| boots | HP=1,MP=1,強度=2,E=5 | HP=2,MP=1,強度=3,E=5 | HP=4,MP=2,強度=4,E=5 | HP=6,MP=5,強度=7,E=5 |

### Diamond（L2〜5）
| slot | L2 | L3 | L4 | L5 |
|------|----|----|----|----|
| helmet | HP=0,MP=2,強度=2,E=- | HP=0,MP=4,強度=3,E=- | HP=2,MP=7,強度=4,E=- | HP=4,MP=15,強度=6,E=- |
| chestplate | HP=0,MP=7,強度=3,E=- | HP=0,MP=13,強度=5,E=6 | HP=3,MP=20,強度=6,E=6 | HP=6,MP=30,強度=14,E=6 |
| leggings | HP=0,MP=7,強度=3,E=- | HP=0,MP=13,強度=5,E=6 | HP=3,MP=20,強度=6,E=6 | HP=6,MP=30,強度=14,E=6 |
| boots | HP=0,MP=2,強度=2,E=- | HP=0,MP=4,強度=3,E=- | HP=2,MP=7,強度=4,E=- | HP=4,MP=15,強度=6,E=- |

### Iron（L1〜5）
| slot | L1 | L2 | L3 | L4 | L5 |
|------|----|----|----|----|-----|
| helmet | HP=0,MP=2,強度=0,E=-,atk+1 | HP=0,MP=2,強度=2,E=-,atk+1 | HP=0,MP=2,強度=3,E=-,atk+1 | HP=2,MP=3,強度=4,E=-,atk+1 | HP=4,MP=10,強度=6,E=-,atk+2 |
| chestplate | HP=0,MP=3,強度=1,E=-,atk+1 | HP=0,MP=4,強度=3,E=-,atk+2 | HP=0,MP=8,強度=4,E=6,atk+3 | HP=3,MP=12,強度=5,E=6,atk+3 | HP=6,MP=15,強度=14,E=6,atk+5 |
| leggings | HP=0,MP=3,強度=1,E=-,atk+1 | HP=0,MP=4,強度=3,E=-,atk+2 | HP=0,MP=8,強度=4,E=6,atk+3 | HP=3,MP=12,強度=5,E=6,atk+3 | HP=6,MP=15,強度=14,E=6,atk+5 |
| boots | HP=0,MP=2,強度=0,E=-,atk+1 | HP=0,MP=2,強度=2,E=-,atk+1 | HP=0,MP=2,強度=3,E=-,atk+1 | HP=2,MP=3,強度=4,E=-,atk+1 | HP=4,MP=10,強度=6,E=-,atk+2 |

### Chainmail（L1〜5）
| slot | L1 | L2 | L3 | L4 | L5 |
|------|----|----|----|----|-----|
| helmet | HP=0,MP=2,強度=0,E=- | HP=0,MP=2,強度=1,E=- | HP=0,MP=4,強度=2,E=- | HP=2,MP=5,強度=2,E=5 | HP=3,MP=10,強度=5,E=- |
| chestplate | HP=0,MP=4,強度=1,E=- | HP=0,MP=7,強度=2,E=- | HP=0,MP=10,強度=2,E=6 | HP=2,MP=18,強度=6,E=5 | HP=4,MP=30,強度=11,E=6 |
| leggings | HP=0,MP=4,強度=1,E=- | HP=0,MP=7,強度=2,E=- | HP=0,MP=10,強度=2,E=6 | HP=2,MP=18,強度=6,E=5 | HP=4,MP=30,強度=11,E=6 |
| boots | HP=0,MP=2,強度=0,E=- | HP=0,MP=2,強度=1,E=- | HP=0,MP=4,強度=2,E=- | HP=2,MP=5,強度=2,E=5 | HP=3,MP=10,強度=5,E=- |

### Copper（L1〜5）
| slot | L1 | L2 | L3 | L4 | L5 |
|------|----|----|----|----|-----|
| helmet | MP=2 | HP=0,MP=2,強度=0,E=0 | HP=0,MP=4,強度=0,E=- | HP=1,MP=6,強度=1,E=5 | HP=3,MP=10,強度=4,E=5 |
| chestplate | MP=2 | HP=0,MP=3,強度=1,E=5 | HP=0,MP=8,強度=2,E=5 | HP=1,MP=10,強度=5,E=5 | HP=5,MP=20,強度=7,E=5 |
| leggings | MP=2 | HP=0,MP=3,強度=1,E=5 | HP=0,MP=8,強度=2,E=5 | HP=1,MP=10,強度=5,E=5 | HP=5,MP=20,強度=7,E=5 |
| boots | MP=2 | HP=0,MP=2,強度=0,E=0 | HP=0,MP=4,強度=0,E=- | HP=1,MP=6,強度=1,E=5 | HP=3,MP=10,強度=4,E=5 |

### Gold（L1〜5）
| slot | L1 | L2 | L3 | L4 | L5 |
|------|----|----|----|----|-----|
| helmet | MP=4 | HP=0,MP=7,強度=0,E=0 | HP=0,MP=7,強度=0,E=- | HP=1,MP=10,強度=1,E=5 | HP=3,MP=20,強度=4,E=5 |
| chestplate | MP=7 | HP=0,MP=10,強度=1,E=6 | HP=0,MP=20,強度=2,E=6 | HP=1,MP=32,強度=5,E=5 | HP=5,MP=50,強度=7,E=5 |
| leggings | MP=7 | HP=0,MP=10,強度=1,E=6 | HP=0,MP=20,強度=2,E=6 | HP=1,MP=32,強度=5,E=5 | HP=5,MP=50,強度=7,E=5 |
| boots | MP=4 | HP=0,MP=7,強度=0,E=0 | HP=0,MP=7,強度=0,E=- | HP=1,MP=10,強度=1,E=5 | HP=3,MP=20,強度=4,E=5 |

### Leather（L1〜5）
| slot | L1 | L2 | L3 | L4 | L5 |
|------|----|----|----|----|-----|
| helmet | MP=5 | MP=8,強度=0,E=0 | HP=0,MP=10,強度=0,E=0 | HP=0,MP=20,強度=1,E=5 | HP=2,MP=30,強度=3,E=5 |
| chestplate | MP=10 | MP=15,強度=1,E=0 | HP=0,MP=30,強度=2,E=5 | HP=0,MP=40,強度=3,E=5 | HP=3,MP=70,強度=5,E=5 |
| leggings | MP=10 | MP=15,強度=1,E=0 | HP=0,MP=30,強度=2,E=5 | HP=0,MP=40,強度=3,E=5 | HP=3,MP=70,強度=5,E=5 |
| boots | MP=5 | MP=8,強度=0,E=0 | HP=0,MP=10,強度=0,E=0 | HP=0,MP=20,強度=1,E=5 | HP=2,MP=30,強度=3,E=5 |

## 生成する JSON エントリ構造

```json
{
    "id": "level{N}_{material}_{slot}",
    "maf": { "maxmp": <MP> },         // MP=0/ブランクなら maf ごと省略
    "minecraft": {
        "components": {
            "minecraft:attribute_modifiers": [
                // 常時: vanilla armor 値
                {"id": "level{N}-{material}-{slot}-armor", "type": "armor", "amount": <vanilla>, "operation": "add_value", "slot": "<mc_slot>"},
                // HP > 0 の場合
                {"id": "level{N}-{material}-{slot}-max-health", "type": "max_health", "amount": <HP>, "operation": "add_value", "slot": "<mc_slot>"},
                // 強度 > 0 の場合
                {"id": "level{N}-{material}-{slot}-armor-toughness", "type": "armor_toughness", "amount": <強度>, "operation": "add_value", "slot": "<mc_slot>"},
                // attribute="attack_damage+N" の場合
                {"id": "level{N}-{material}-{slot}-attack-damage", "type": "attack_damage", "amount": <N>, "operation": "add_value", "slot": "<mc_slot>"}
            ],
            // Enchant > 0 の場合
            "minecraft:enchantments": {
                "minecraft:protection": <Enchant>
            }
        }
    },
    "itemId": "minecraft:{material}_{slot}"
}
```

スロット名変換: helmet→head, chestplate→chest, leggings→legs, boots→feet

## 実装手順

1. 一時的な Python スクリプト `/tmp/gen_armor.py` を書く
   - 上記全スペックをデータとして持つ
   - `savedata/item/armors/level{N}.json` を生成
2. スクリプト実行
3. スクリプト削除
4. `cd maf-command-generator && make check` で検証
5. エラーがあれば修正して再実行

## 検証

- `make check` でビルドとテストが通ること
- level1〜5.json が `savedata/item/armors/` に存在すること
- 各ファイルが valid JSON であること（`make check` のバリデーションで確認）
- copper の itemId (`minecraft:copper_helmet` 等) がバリデーションで通ること

## 注意点

- copper armor は vanilla 1.21 以前には存在しない。バリデーションエラーが出た場合はアイテムIDを確認して修正
- E=0 は「Enchant値0」= エンチャントなし（E=-と同扱い）でスキップ
- HP=0 もスキップ（max_health=0 のモディファイアは不要）
- 強度=0 もスキップ（armor_toughness=0 のモディファイアは不要）
