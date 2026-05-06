---
name: add-maf-enemy-skill
description: maf-command-generatorのsavedataにエネミースキル（敵モブが使用するスキル）を追加するスキル。enemy_skill, エネミースキル, 敵スキル, モブスキル, 毒攻撃, 範囲攻撃 などの新規追加が出たら使う。savedata/enemy_skill/ai_workspace.json への書き込みを担う。
disable-model-invocation: true
---

# エネミースキル追加

**対象ファイル:** `maf-command-generator/savedata/enemy_skill/ai_workspace.json`

エネミースキルはモブが一定間隔で発動するスキル。Enemyエンティティの `enemySkillIds` から参照される。

**実装方針:** `savedata/` 配下のJSONのみ参照・編集する。`maf-command-generator/app/` の Goコードは読まない。

## 要件整理

実装前に AskUserQuestion で不明点を確認する。

- スキルID（例: near_poison ／ 小文字・アンダースコア・ハイフンのみ）
- スキル説明
- 効果の詳細（範囲・対象・強さ・演出）
- このスキルを使わせるエネミーID（`enemySkillIds` への紐付け要否）

## 手順

1. `ai_workspace.json` が存在すれば読み込む。なければ `{"entries": []}` として扱う
2. `entries` に新エントリを追加
3. ファイルに書き戻す
4. `cd maf-command-generator && make run/export` でエクスポート
5. `make mc-cmd CMD='reload'` でデータパックリロード
6. `make mc-cmd CMD='function maf:generated/enemy/skill/[id]'` を実行し、"Unknown or incomplete command" が出ないことを確認（文法チェックのみ・動作の正確さは検証外）

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

作成したエネミースキルは、`add-maf-enemy` スキル（またはenemyのai_workspace.json）で `enemySkillIds` にIDを追加して使用する:

```json
"enemySkillIds": ["near_poison", "my_new_skill"]
```

## 参照スキル

| スキル | 用途 |
|--------|------|
| `enemy-replace` | エネミーのスポーンシステム・発動タイミング詳細 |
| `rcon` | RCONコマンド発行方法（make mc-cmd の使い方） |
