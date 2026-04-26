---
name: add-maf-passive
description: maf-command-generatorのsavedataにパッシブスキルを追加するスキル。passive, パッシブ, 常時発動, 攻撃時発動, スロット装備スキル などの新規追加が出たら使う。savedata/passive/claude.json への書き込みを担う。
disable-model-invocation: true
---

# パッシブスキル追加

**対象ファイル:** `maf-command-generator/savedata/passive/claude.json`

**実装方針:** `savedata/` 配下のJSONのみ参照・編集する。`maf-command-generator/app/` の Goコードは読まない。

## 要件整理

実装前に AskUserQuestion で不明点を確認する。

- パッシブID（例: regeneration ／ 小文字・アンダースコア・ハイフンのみ）
- パッシブ名（表示名）
- 発動条件（always / attack / none）
- 装備スロット（1〜3、複数可）
- 効果の詳細
- グリモア生成要否（ドロップから入手させるか）

## 手順

1. `claude.json` が存在すれば読み込む。なければ `{"entries": []}` として扱う
2. `entries` に新エントリを追加
3. ファイルに書き戻す
4. `cd maf-command-generator && make run/export` でエクスポート
5. `make mc-cmd CMD='reload'` でデータパックリロード
6. `make mc-cmd CMD='function maf:generated/passive/give/[id]_slot[N]'` を実行し、"Unknown or incomplete command" が出ないことを確認（slot番号は定義したslots[0]を使う。文法チェックのみ・動作の正確さは検証外）

## JSONスキーマ

```json
{
  "entries": [
    {
      "id": "passive_id",          // 必須: 小文字・アンダースコア・ハイフンのみ、ユニーク
      "name": "パッシブ名",         // 必須: 表示名(80文字以内)
      "role": "役割説明",           // 必須: 役割の説明(200文字以内)
      "condition": "always",       // 必須: "always" | "attack" | "none"
      "slots": [1],                // 必須: 装備可能スロット番号(1〜3の配列、1つ以上)
      "script": [                  // 必須: 発動時コマンド(1行以上)
        "effect give @s minecraft:regeneration 10 0"
      ],
      "generate_grimoire": false   // 必須: true=設定書生成 / false=設定書なし
    }
  ]
}
```

## conditionの意味

| condition | 発動タイミング |
|-----------|-------------|
| `always` | 毎tick（装備中常時） |
| `attack` | 近接攻撃でダメージを与えたとき |
| `none` | 自動発動なし（他functionから手動呼び出し） |

## generate_grimoireの意味

| 値 | 設定書(give/apply)生成 | ルートテーブルから入手 |
|---|---|---|
| `true` | される | できる |
| `false` | されない | できない |

- ドロップやトレジャーから入手可能にしたいなら `true`
- 特定アイテムにのみ付与する場合は `false`

## scriptのパターン

```mcfunction
# always条件: バフ維持
execute unless entity @s[nbt={active_effects:[{id:"minecraft:regeneration"}]}] run effect give @s regeneration 10 0

# attack条件: 攻撃時エフェクト
effect give @e[distance=..4,type=!minecraft:player] minecraft:poison 5 1

# プレイヤー中心の演出
particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 3
```

## slots指定例

- `[1]` → スロット1のみ装備可
- `[1, 2, 3]` → どのスロットにも装備可
- `[2, 3]` → スロット2か3のみ

## 参照スキル

| スキル | 用途 |
|--------|------|
| `passive` | パッシブのデータモデル・発動条件・スロットシステム詳細 |
| `magic-casting` | generate_grimoire=trueのとき詠唱パイプラインも適用される |
| `rcon` | RCONコマンド発行方法（make mc-cmd の使い方） |
