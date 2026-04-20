---
name: add-maf-grimoire
description: maf-command-generatorのsavedataにグリモア（魔法/スペル）を追加するスキル。grimoire, 魔法, スペル, 呪文, 魔導書 などの追加・新規作成が出たら使う。savedata/grimoire/claude.json への書き込みを担う。
disable-model-invocation: true
---

# グリモア追加

**対象ファイル:** `maf-command-generator/savedata/grimoire/claude.json`

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
      "id": "spell_name01",        // 必須: 小文字・アンダースコア・ハイフンのみ、ユニーク
      "castTime": 40,              // 必須: 詠唱時間(tick)。通常40(2秒)
      "coolTime": 20,              // 必須: クールダウン(tick)。通常20(1秒)
      "mpCost": 10,                // 必須: 消費MP(0〜60程度)
      "script": [                  // 必須: 発動時mcfunctionコマンド(1行以上)
        "effect give @e[distance=..8,type=!#maf:undead] minecraft:instant_damage 1 1",
        "playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2",
        "tellraw @a[distance=..50] [{\"selector\":\"@s\"},{\"text\":\" は スペル名 を唱えた！\"}]"
      ],
      "title": "スペル名",         // 必須: 表示名
      "description": "効果説明"    // 任意
    }
  ]
}
```

## scriptの典型パターン

```mcfunction
# ダメージ系(undead逆効果に注意)
execute as @e[distance=1..8,type=#maf:undead] run effect give @s minecraft:instant_health 1 1
execute as @e[distance=1..8,type=!#maf:undead] run effect give @s minecraft:instant_damage 1 1

# 味方バフ
effect give @e[type=#maf:friendmob,distance=..10] minecraft:regeneration 10 0

# 演出
playsound minecraft:entity.blaze.shoot master @a ~ ~ ~ 1.0 0.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は スペル名 を唱えた！"}]
```

## エンティティタグ

| タグ | 意味 |
|------|------|
| `#maf:undead` | アンデッド（ゾンビ・スケルトン等） |
| `#maf:enemymob` | 敵モブ全般 |
| `#maf:friendmob` | 味方モブ |

## mpCost目安

| 強さ | mpCost |
|------|--------|
| 弱い | 4〜8 |
| 中程度 | 9〜18 |
| 強い | 19〜30 |
| 超強力 | 31〜60 |
