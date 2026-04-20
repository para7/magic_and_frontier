---
name: add-maf-enemy-skill
description: maf-command-generatorのsavedataにエネミースキル（敵モブが使用するスキル）を追加するスキル。enemy_skill, エネミースキル, 敵スキル, モブスキル, 毒攻撃, 範囲攻撃 などの新規追加が出たら使う。savedata/enemy_skill/claude.json への書き込みを担う。
disable-model-invocation: true
---

# エネミースキル追加

**対象ファイル:** `maf-command-generator/savedata/enemy_skill/claude.json`

エネミースキルはモブが一定間隔で発動するスキル。Enemyエンティティの `enemySkillIds` から参照される。

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
      "id": "skill_id",            // 必須: 小文字・アンダースコア・ハイフンのみ、ユニーク
      "description": "スキル説明",  // 必須: スキルの説明
      "script": [                  // 必須: 発動時mcfunctionコマンド(1行以上)
        "effect give @e[distance=..3] minecraft:poison 10 2"
      ]
    }
  ]
}
```

## scriptのパターン

エネミースキルは **モブ自身が実行者(@s)** として呼ばれる。座標基点はモブの位置。

```mcfunction
# 範囲毒
effect give @e[distance=..4] minecraft:poison 10 2

# 近距離スロー
effect give @e[distance=..3,type=minecraft:player] minecraft:slowness 5 3

# 爆発系
summon creeper ~ ~ ~ {NoAI:1b,powered:1b,Fuse:0,ExplosionRadius:2b}

# 召喚
summon minecraft:zombie ~ ~ ~ {CustomName:'{"text":"雑魚ゾンビ"}'}

# 演出
particle minecraft:flame ~ ~1 ~ 0.5 0.5 0.5 0.1 30
playsound minecraft:entity.blaze.shoot master @a ~ ~ ~ 1 0.8
```

## Enemyへの紐付け

作成したエネミースキルは、`add-maf-enemy` スキル（またはenemyのclaude.json）で `enemySkillIds` にIDを追加して使用する:

```json
"enemySkillIds": ["near_poison", "my_new_skill"]
```
