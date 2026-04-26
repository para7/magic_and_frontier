---
name: add-maf-enemy
description: maf-command-generatorのsavedataにカスタムエネミー（敵モブ）を追加するスキル。enemy, エネミー, 敵, カスタムモブ, ボス, ゾンビ強化, モブ追加 などの新規作成が出たら使う。savedata/enemy/claude.json への書き込みを担う。
disable-model-invocation: true
---

# カスタムエネミー追加

**対象ファイル:** `maf-command-generator/savedata/enemy/claude.json`

**実装方針:** `savedata/` 配下のJSONのみ参照・編集する。`maf-command-generator/app/` の Goコードは読まない。

## 要件整理

実装前に AskUserQuestion で不明点を確認する。

- エネミーID（例: poison_zombie ／ 小文字・アンダースコア・ハイフンのみ）
- 名前・ベースモブ
- HP倍率・攻撃力・防御力・移動速度
- 使用するエネミースキルID（任意）
- ドロップ内容（アイテムID・グリモアID・重み・個数）
- dropMode（replace / append）

## 手順

1. `claude.json` が存在すれば読み込む。なければ `{"entries": []}` として扱う
2. `entries` に新エントリを追加
3. ファイルに書き戻す
4. `cd maf-command-generator && make run/export` でエクスポート
5. `make mc-cmd CMD='reload'` でデータパックリロード
6. `make mc-cmd CMD='function maf:generated/enemy/spawn/[id]'` を実行し、"Unknown or incomplete command" が出ないことを確認（文法チェックのみ・動作の正確さは検証外）

## JSONスキーマ

```json
{
  "entries": [
    {
      "id": "enemy_id",            // 必須: 小文字・アンダースコア・ハイフンのみ、ユニーク
      "mobType": "minecraft:zombie", // 必須: ベースとなるバニラモブID
      "name": "エネミー名",          // 必須: 表示名
      "hp": 2,                     // 必須: HP倍率(バニラ基準の倍率、整数)
      "memo": "",                  // 任意: メモ
      "attack": 6,                 // 必須: 攻撃力
      "defense": 2,                // 必須: 防御力
      "moveSpeed": 0.22,           // 必須: 移動速度(バニラゾンビ=0.23)
      "equipment": {},             // 必須: 装備(通常空オブジェクト)
      "enemySkillIds": [],         // 必須: 使用するenemySkillのIDリスト
      "dropMode": "replace",       // 必須: "replace" | "append"
      "drops": [                   // 必須: ドロップテーブル
        {
          "rolls": 1,
          "entries": [
            {
              "type": "maf:item",      // "maf:item" | "maf:grimoire"
              "name": "item_id",       // アイテムID or グリモアID
              "weight": 70,            // 抽選重み
              "count": 1               // 個数(整数 or {"min":1,"max":3})
            }
          ]
        }
      ]
    }
  ]
}
```

## dropModeの意味

| 値 | 意味 |
|----|------|
| `replace` | バニラのドロップを完全に置き換える |
| `append` | バニラのドロップに追加する |

## dropsの型バリエーション

```json
// 固定個数
{"type": "maf:item", "name": "sword_01", "weight": 80, "count": 1}

// ランダム個数
{"type": "maf:item", "name": "herb_01", "weight": 50, "count": {"min": 1, "max": 3}}

// グリモアドロップ
{"type": "maf:grimoire", "name": "healing01", "weight": 20, "count": 1}
```

## よく使うmobType

- `minecraft:zombie`, `minecraft:skeleton`, `minecraft:spider`
- `minecraft:cave_spider`, `minecraft:creeper`, `minecraft:witch`
- `minecraft:enderman`, `minecraft:blaze`, `minecraft:wither_skeleton`

## 参照スキル

| スキル | 用途 |
|--------|------|
| `enemy-replace` | スポーンテーブルシステム・置換ロジック詳細 |
| `rcon` | RCONコマンド発行方法（make mc-cmd の使い方） |
