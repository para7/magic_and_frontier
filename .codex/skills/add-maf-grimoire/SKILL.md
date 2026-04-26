---
name: add-maf-grimoire
description: maf-command-generatorのsavedataにグリモア（魔法/スペル）を追加するスキル。grimoire, 魔法, スペル, 呪文, 魔導書 などの追加・新規作成が出たら使う。savedata/grimoire/claude.json への書き込みを担う。
disable-model-invocation: true
---

# グリモア追加

**対象ファイル:** `maf-command-generator/savedata/grimoire/claude.json`

ユーザー生成物とAI生成物を分けるため、claude.json に必ず出力すること。ファイルがなければ新規作成する。

**実装方針:** `savedata/` 配下のJSONのみ参照・編集する。`maf-command-generator/app/` の Goコードは読まない。

## 要件整理

実装前に AskUserQuestion で不明点を確認する。

- グリモアID（例: thunder01 ／ 小文字・アンダースコア・ハイフンのみ）
- スペル名（表示名）
- 効果（範囲・対象・強さ）
- castTime / coolTime / mpCost（未定なら目安を提案して確認）
- 演出（パーティクル・サウンド）の有無・種類

## 手順

1. `claude.json` が存在すれば読み込む。なければ `{"entries": []}` として扱う
2. `entries` に新エントリを追加
3. ファイルに書き戻す
4. `cd maf-command-generator && make run/export` でエクスポート
5. `make mc-cmd CMD='reload'` でデータパックリロード
6. `make mc-cmd CMD='function maf:generated/grimoire/give/[id]'` を実行し、"Unknown or incomplete command" が出ないことを確認（文法チェックのみ・動作の正確さは検証外）

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

## 参照スキル

| スキル | 用途 |
|--------|------|
| `grimoire` | グリモアNBT構造・ランタイムフロー・set_count等の詳細 |
| `magic-casting` | castTime/coolTime/mpCostの制約・詠唱パイプライン全体像 |
| `rcon` | RCONコマンド発行方法（make mc-cmd の使い方） |
