---
name: add-maf-bow
description: maf-command-generatorのsavedataに弓パッシブ（矢の効果スキル）を追加するスキル。bow, 弓, 矢, 弓パッシブ, 矢効果, 着弾エフェクト, 飛行エフェクト などの新規追加が出たら使う。savedata/bow/ai_workspace.json への書き込みを担う。
disable-model-invocation: true
---

# 弓パッシブ追加

**対象ファイル:** `maf-command-generator/savedata/bow/ai_workspace.json` ユーザー生成物とAI生成物を分けるため、ai_workspace.json に必ず出力すること。ファイルがなければ新規作成する。

弓パッシブは矢の挙動（発射・飛行・着弾・命中）に付与するスキル。Passiveとは別エンティティ。

**実装方針:** `savedata/` 配下のJSONのみ参照・編集する。`maf-command-generator/app/` の Goコードは読まない。

## 要件整理

実装前に AskUserQuestion で不明点を確認する。

- 弓パッシブID（例: poison_arrow ／ 小文字・アンダースコア・ハイフンのみ）
- 名前・役割説明
- 発動フェーズ（hit / fired / flying / ground ／ 複数可）
- 各フェーズの効果詳細

## 手順

1. `ai_workspace.json` が存在すれば読み込む。なければ `{"entries": []}` として扱う
2. `entries` に新エントリを追加
3. ファイルに書き戻す
4. `cd maf-command-generator && make run/export` でエクスポート
5. `make mc-cmd CMD='reload'` でデータパックリロード
6. 定義したscriptフェーズに対応するfunctionを実行して確認（例: flying定義なら `function maf:generated/bow/flying/[id]_flying`）。"Unknown or incomplete command" が出ないことを確認（文法チェックのみ・動作の正確さは検証外）

## JSONスキーマ

```json
{
  "entries": [
    {
      "id": "bow_skill_id",        // 必須: 小文字・アンダースコア・ハイフンのみ、ユニーク
      "name": "弓スキル名",         // 必須: 表示名
      "role": "役割説明",           // 必須: 説明
      "slots": [],                 // 必須: 通常空配列
      "life_sub": 100,             // 必須: 矢の耐久消費量(通常100)
      "script_hit": ["..."],       // 任意: 矢がエンティティに命中したとき
      "script_fired": ["..."],     // 任意: 矢を発射したとき(プレイヤー位置)
      "script_flying": ["..."],    // 任意: 矢飛行中(矢の位置、毎tick)
      "script_ground": ["..."]     // 任意: 矢が地面に着弾したとき
    }
  ]
}
```

## スクリプトの使い分け

| フィールド | 実行タイミング | 座標基点 |
|-----------|-------------|---------|
| `script_hit` | 矢がモブ/プレイヤーに当たった瞬間 | 矢の位置 |
| `script_fired` | 弓を放った瞬間 | プレイヤーの位置 |
| `script_flying` | 矢が飛行中(毎tick) | 矢の位置 |
| `script_ground` | 矢が地面に刺さった瞬間 | 矢の位置 |

4つのスクリプトはすべて任意。必要なものだけ定義。少なくとも1つ必要。

## scriptのパターン

```mcfunction
# hit: 命中したエンティティに効果(距離..4が矢の当たり判定範囲)
effect give @e[distance=..4,type=!minecraft:player] minecraft:poison 5 2

# fired: 発射時演出(プレイヤー位置基点)
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 1 1.5 1

# flying: 飛行中パーティクル(矢の位置基点)
particle minecraft:happy_villager ~ ~ ~ 0.4 0.4 0.4 0 5 force

# ground: 着弾時(矢が刺さった場所基点)
summon minecraft:skeleton ~ ~ ~
kill @s
```

## life_subの意味

矢1本の耐久消費量。`100` = 通常の耐久消費（矢が消える）。

## 参照スキル

| スキル | 用途 |
|--------|------|
| `bow-passive` | 弓パッシブの動作フロー・矢へのデータ埋め込み・dolphins_graceマーカー詳細 |
| `rcon` | RCONコマンド発行方法（make mc-cmd の使い方） |
