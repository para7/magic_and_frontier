---
name: bow-passive
description: 弓パッシブの動作フローリファレンス。弓を引く→矢にタグ付け→矢着弾→エフェクト発動の2段階フローを網羅する。矢へのデータ埋め込み、dolphins_grace マーカーによる被弾エンティティ検出、PierceLevel による矢残存など弓パッシブ固有の仕組みを解説。弓パッシブのバグ修正、矢の挙動変更、着弾エフェクトの調整、新規弓パッシブ追加などで参照すること。bow, 弓, arrow, 矢, 着弾, hit, tag_passive_arrow, on_arrow_hit, mafBowUsed などのキーワードで使う。
---

# 弓パッシブ動作フロー

弓パッシブは「弓を引いて矢を放つ → 矢が命中 → 命中先でエフェクト発動」という2段階フローを持つパッシブスキル。通常の `always` パッシブが毎 tick 直接効果を実行するのに対し、弓パッシブは **矢エンティティにスキル情報を載せて運ぶ** 点が根本的に異なる。

## 関連スキル

- **passive**: パッシブシステム全体（データモデル、スロット、設定書）
- **magic-casting**: 詠唱パイプライン（設定書使用時に通る）
- **ohmydat**: プレイヤー個別ストレージ。パッシブ装備状態の保存先
- **maf-export**: Go ジェネレータのエクスポートパイプライン

---

## 1. 全体フロー概要

```
[発射フェーズ — 毎 tick]
passive/tick → run_effect → generated/passive/effect/{id}
  → mafBowUsed >= 1 ? (弓を使ったか)
  → YES → 近くの矢を探す → tag_passive_arrow でデータ埋め込み

[飛翔中]
矢が飛行。PierceLevel:2b により貫通属性あり（命中で消えない）

[着弾検知 — advancement トリガー]
advancement/arrow_hit.json (player_hurt_entity)
  → on_arrow_hit
    → advancement revoke（再発火可能に）
    → mafBowHit += 1（フラグのみ。処理は tick に委譲）

[着弾処理 — 次 tick]
passive/tick → mafBowHit >= 1 → on_bow_hit
  → maf_passive_arrow タグ付き矢を検索
  → 矢の custom_data からパッシブ情報を maf:tmp に退避
  → dolphins_grace(amp:80) を持つエンティティに対して:
     → run_bow_effect → generated/passive/bow/{passiveId}
  → dolphins_grace マーカー除去
  → 半径5以内のタグ付き矢を kill、ストレージクリア
```

---

## 2. 発射フェーズ（矢のタグ付け）

### 2a. トリガー検出

弓パッシブの effect ファイルは Go ジェネレータが自動生成する。

**生成元:** `maf-command-generator/app/domain/export/passive.go` の `buildBowEffectBody()`

**生成先:** `datapacks/.../generated/passive/effect/{id}.mcfunction`

```mcfunction
execute unless score @s mafBowUsed matches 1.. run return 0
execute store result storage maf:tmp bow_player_id int 1 run scoreboard players get @s mafPlayerID
execute as @e[type=arrow,distance=..2,nbt=!{inGround:1b},sort=nearest,limit=1] run function maf:magic/passive/tag_passive_arrow {passive_id:"{id}",life:{lifeValue}}
```

| 行 | 処理 | 詳細 |
|----|------|------|
| 1行目 | 弓使用チェック | `mafBowUsed`（`minecraft.used:minecraft.bow` 統計）が 1 以上でなければ即 return |
| 2行目 | プレイヤーID退避 | `mafPlayerID` を `maf:tmp bow_player_id` にコピー（矢にシューターIDを書き込むため） |
| 3行目 | 矢タグ付け | 半径2以内、地面に刺さっていない最近矢を対象に `tag_passive_arrow` を呼ぶ |

### 2b. mafBowUsed スコアボード

| ファイル | 処理 |
|---------|------|
| `load.mcfunction:11` | `scoreboard objectives add mafBowUsed minecraft.used:minecraft.bow` — バニラ統計から自動カウント |
| `system/score/afterscore.mcfunction:11` | `scoreboard players set @a mafBowUsed 0` — 毎 tick 末にリセット |

弓を使った tick だけ `mafBowUsed >= 1` になる。次 tick にはリセットされるため、1回の射撃で1回だけ矢タグ付けが走る。

### 2c. tag_passive_arrow.mcfunction（矢へのデータ埋め込み）

**ファイル:** `datapacks/.../magic/passive/tag_passive_arrow.mcfunction`

```mcfunction
# @s = arrow entity, マクロ引数: passive_id, life
$data merge entity @s {Tags:["maf_passive_arrow"],PierceLevel:2b,item:{components:{"minecraft:custom_data":{maf:{passiveId:"$(passive_id)"}},"minecraft:potion_contents":{custom_effects:[{id:"minecraft:dolphins_grace",duration:200,amplifier:80}]}}}}
$data modify entity @s life set value $(life)s
data modify entity @s item.components."minecraft:custom_data".maf.shooterPlayerID set from storage maf:tmp bow_player_id
```

矢エンティティに書き込む情報:

| データ | 値 | 目的 |
|--------|-----|------|
| `Tags: ["maf_passive_arrow"]` | 固定 | 着弾時にパッシブ矢を `@e` セレクタで高速検索 |
| `PierceLevel: 2b` | 固定 | 命中しても矢が消えないようにする（着弾処理後に手動 kill） |
| `custom_data.maf.passiveId` | パッシブID | 着弾時にどのパッシブか識別 |
| `potion_contents.custom_effects` | dolphins_grace amp:80 | **被弾マーカー**: 矢が当たったエンティティにこのエフェクトが付く |
| `life` | `1200 - LifeSub` | 矢の生存 tick（短いほど早く消滅、0 で即消滅） |
| `custom_data.maf.shooterPlayerID` | プレイヤーID | 将来用。どのプレイヤーが撃ったか |

### 2d. life 値の計算

```go
// passive.go
lifeSub := 1200  // デフォルト: 即消滅
if entry.Bow != nil && entry.Bow.LifeSub != nil {
    lifeSub = *entry.Bow.LifeSub
}
lifeValue := 1200 - lifeSub
```

| LifeSub | life 値 | 矢の飛行時間 |
|---------|---------|-------------|
| 未設定 (nil) | 0 | 即消滅（着弾前に消える可能性大） |
| 0 | 1200 | 60秒（バニラ標準） |
| 100 | 1100 | 55秒 |
| 1200 | 0 | 即消滅 |

> **注意:** LifeSub 未設定時のデフォルトは 1200（life=0、即消滅）。実用的な弓パッシブでは必ず `LifeSub` を明示的に指定すること。

---

## 3. 着弾フェーズ（エフェクト発動）

着弾処理は **2段階**に分離されている:
1. advancement が即座にフラグを立てる（`on_arrow_hit`）
2. 次の tick ループで実際の処理を行う（`on_bow_hit`）

### 3a. Advancement トリガー

**ファイル:** `datapacks/.../advancement/arrow_hit.json`

```json
{
  "criteria": {
    "arrow_hit": {
      "trigger": "minecraft:player_hurt_entity",
      "conditions": {
        "damage": {
          "type": {
            "direct_entity": { "type": "minecraft:arrow" }
          }
        }
      }
    }
  },
  "rewards": {
    "function": "maf:magic/passive/on_arrow_hit"
  }
}
```

- `player_hurt_entity` — **プレイヤーが矢でエンティティを傷つけた** ときに発火
- `@s` は **攻撃側プレイヤー**（被弾エンティティではない）

### 3b. on_arrow_hit.mcfunction（フラグセットのみ）

**ファイル:** `datapacks/.../magic/passive/on_arrow_hit.mcfunction`

```mcfunction
advancement revoke @s only maf:arrow_hit
scoreboard players add @s mafBowHit 1
```

advancement のコールバックでは**フラグを立てるだけ**。実際の処理は tick に委譲する。`mafBowHit` は dummy スコアボードで、`load.mcfunction` で登録される。

### 3c. mafBowHit スコアボード

| 項目 | 内容 |
|------|------|
| 登録 | `load.mcfunction` — `scoreboard objectives add mafBowHit dummy` |
| 加算 | `on_arrow_hit.mcfunction` — `scoreboard players add @s mafBowHit 1` |
| リセット | `on_bow_hit.mcfunction` — tick 処理の先頭で `set @s mafBowHit 0` |
| `afterscore` でのリセット | **しない** — tick 処理まで値を保持する必要があるため |

### 3d. passive/tick.mcfunction（フラグチェック）

```mcfunction
# 弓着弾処理
execute if score @s mafBowHit matches 1.. run function maf:magic/passive/on_bow_hit
```

`passive/tick` の末尾に追加。`magic/tick` から `execute as @a at @s` で呼ばれるため、プレイヤーごとに評価される。

### 3e. on_bow_hit.mcfunction（着弾処理本体）

**ファイル:** `datapacks/.../magic/passive/on_bow_hit.mcfunction`

```mcfunction
# フラグをリセット
scoreboard players set @s mafBowHit 0                             # ① フラグクリア

# タグ付き矢が近くになければ終了
execute unless entity @e[type=arrow,tag=maf_passive_arrow,        # ② パッシブ矢がなければ終了
  limit=1,sort=nearest,distance=..10] run return 0

# 矢のパッシブ情報を退避（kill 前に取得）
data modify storage maf:tmp arrow_passive set from entity         # ③ custom_data.maf を退避
  @e[type=arrow,tag=maf_passive_arrow,sort=nearest,limit=1]
  item.components."minecraft:custom_data".maf

# dolphins_grace(amp:80) を持つエンティティに対してエフェクト発動
execute as @e[nbt={active_effects:[                               # ④ 被弾エンティティでエフェクト実行
  {id:"minecraft:dolphins_grace",amplifier:80b}]}] at @s run function
  maf:magic/passive/run_bow_effect with storage maf:tmp arrow_passive

# マーカーエフェクトを除去
effect clear @e[nbt={active_effects:[                             # ⑤ dolphins_grace マーカー除去
  {id:"minecraft:dolphins_grace",amplifier:80b}]}]
  minecraft:dolphins_grace

# 半径5以内の最近タグ付き矢を kill
kill @e[type=arrow,tag=maf_passive_arrow,sort=nearest,limit=1,    # ⑥ 矢を手動 kill
  distance=..5]

data remove storage maf:tmp arrow_passive                         # ⑦ ストレージクリア
```

**処理の流れ:**

| ステップ | 処理 | なぜ必要か |
|---------|------|----------|
| ① フラグクリア | `mafBowHit` を 0 にリセット | 次 tick に再処理されないよう先にクリア |
| ② 矢存在チェック | タグ付き矢が 10ブロック以内にあるか | 通常の矢（パッシブなし）を無視 |
| ③ データ退避 | 矢の `custom_data.maf` を `maf:tmp arrow_passive` にコピー | kill 後はデータ取得不可なので先に退避 |
| ④ エフェクト実行 | dolphins_grace (amp:80) を持つエンティティを `@s` として実行 | 矢の命中で効果が付与された被弾エンティティを正確に特定 |
| ⑤ マーカー除去 | dolphins_grace (amp:80) を除去 | 副作用（泳ぎ速度上昇）の防止 |
| ⑥ 矢 kill | 半径5以内のタグ付き矢を除去 | `PierceLevel` で残存した矢の後始末 |
| ⑦ クリア | 一時ストレージ削除 | ゴミデータ残留防止 |

### 3f. 被弾エンティティの検出方式（dolphins_grace マーカー方式）

```mcfunction
execute as @e[nbt={active_effects:[{id:"minecraft:dolphins_grace",amplifier:80b}]}] at @s run ...
```

矢の `potion_contents` に仕込んだ dolphins_grace (amplifier:80) が命中時に被弾エンティティへ付与される。このエフェクトを持つエンティティ＝今の矢で被弾したエンティティとして正確に特定できる。

| 特性 | 内容 |
|------|------|
| 特定精度 | 高い。amplifier:80 は通常付与されない特殊値なので誤検出が起きにくい |
| 複数命中 | 矢が貫通して複数エンティティに当たった場合も、全員がマーカーを持つため全員にエフェクト発動 |
| 無敵時間 | 無敵時間中でもエフェクト付与は起きるため、HurtTime 方式より検出漏れが少ない |
| クリア必須 | エフェクトが残ると泳ぎ速度上昇の副作用が出るため、エフェクト実行後に必ず `effect clear` する |

`tag_passive_arrow.mcfunction` での矢への埋め込み:
```mcfunction
"minecraft:potion_contents":{custom_effects:[{id:"minecraft:dolphins_grace",duration:200,amplifier:80}]}
```
duration:200 tick（10秒）で自然消滅するが、on_bow_hit 内で即クリアするため通常は残らない。

### 3g. run_bow_effect.mcfunction（ディスパッチ）

**ファイル:** `datapacks/.../magic/passive/run_bow_effect.mcfunction`

```mcfunction
$function maf:generated/passive/bow/$(passiveId)
```

`maf:tmp arrow_passive` ストレージから `passiveId` を受け取り、マクロで `generated/passive/bow/{id}.mcfunction` にディスパッチする。

### 3e. 生成された bow スクリプト

**ファイル:** `datapacks/.../generated/passive/bow/{id}.mcfunction`

`savedata/passive.json` の `Script[]` がそのまま出力される。**`@s` は被弾エンティティ**、座標も被弾エンティティの位置。

例（`test_bow_passive`）:
```mcfunction
effect give @e[distance=..4] minecraft:glowing 3 10
effect give @e[distance=..4,type=!minecraft:player] minecraft:levitation 1 2
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 1 2 1
tell @a "矢の効果が発動しました"
summon minecraft:zombie ~ ~ ~
```

---

## 4. 生成パイプライン（Go 側）

### 4a. 成果物の分離

弓パッシブは **2種類の mcfunction** が生成される点が `always` と異なる:

| 成果物 | パス | 内容 | 記述者 |
|--------|------|------|--------|
| effect ファイル | `generated/passive/effect/{id}.mcfunction` | 弓検知 + 矢タグ付け | **自動生成**（`buildBowEffectBody()`） |
| bow ファイル | `generated/passive/bow/{id}.mcfunction` | 着弾時エフェクト | **ユーザー記述**（`Script[]`） |

`always` パッシブでは `Script[]` が effect ファイルに直接出力されるが、弓パッシブでは effect ファイルは自動生成コードで埋められ、`Script[]` は bow ファイルに分離される。

### 4b. BuildPassiveArtifacts の弓パッシブ分岐

**ファイル:** `maf-command-generator/app/domain/export/passive.go`

```go
// condition == "bow" の場合
lifeSub := 1200
if entry.Bow != nil && entry.Bow.LifeSub != nil {
    lifeSub = *entry.Bow.LifeSub
}

// effect ファイル → 自動生成（弓検知 + 矢タグ付け）
effects = append(effects, PassiveEffectFunction{
    ID:   entry.ID,
    Body: buildBowEffectBody(entry.ID, lifeSub),
})

// bow ファイル → ユーザーの Script[]
bows = append(bows, PassiveBowFunction{
    ID:   entry.ID,
    Body: strings.Join(entry.Script, "\n"),
})
```

### 4c. ファイル書き出し

**ファイル:** `maf-command-generator/app/domain/export/passive.go` の `WritePassiveArtifacts()`

```
WritePassiveArtifacts()
  ├── effect/{id}.mcfunction   ← 全パッシブ（always: Script, bow: 自動生成）
  ├── bow/{id}.mcfunction      ← bow パッシブのみ（Script[]）
  ├── give/{id}_slot{N}.mcfunction  ← 全パッシブ
  └── apply/{id}_slot{N}.mcfunction ← 全パッシブ
```

---

## 5. データの流れ（矢を介した情報伝達）

```
[発射時]
oh_my_dat passive.slot{N}.id → run_effect → effect/{id}
  → maf:tmp bow_player_id ← mafPlayerID
  → 矢エンティティに書き込み:
     ├── Tags: ["maf_passive_arrow"]
     ├── custom_data.maf.passiveId: "{id}"
     ├── custom_data.maf.shooterPlayerID: {playerID}
     ├── potion_contents: dolphins_grace(amp:80)
     └── life: {1200 - LifeSub}

[着弾時]
矢 custom_data.maf → maf:tmp arrow_passive
  → run_bow_effect with {passiveId: "{id}"}
    → generated/passive/bow/{id}.mcfunction (@s = 被弾エンティティ)
```

---

## 6. 設計上の注意点

### PierceLevel の必要性

通常の矢は命中で消滅する。矢が消えると `custom_data` が失われるため、`PierceLevel:2b` で貫通属性を付けて矢を残存させる。着弾処理後に手動で `kill` する。

### dolphins_grace マーカー方式の注意点

- 矢が複数のエンティティを貫通した場合、全員にエフェクトが発動する（PierceLevel:2b のため）
- `effect clear` を忘れるとエンティティに泳ぎ速度上昇が残り続ける
- duration:200 tick の自然消滅前に必ずクリアすること

### mafBowUsed のリセットタイミング

`afterscore` は tick 末で実行。passive/tick は magic/tick から呼ばれるため、弓使用 → 同一 tick 内で矢タグ付け → tick 末にリセット、という流れ。

### LifeSub デフォルト値の罠

`BowConfig.LifeSub` が nil（未設定）の場合、デフォルトの lifeSub は **1200** になり、`life = 1200 - 1200 = 0` で矢が即消滅する。新規弓パッシブには必ず `LifeSub` を明示的に設定すること。

---

## 7. 弓パッシブ追加手順

1. `savedata/passive.json` にエントリ追加
   - `condition: "bow"` を設定
   - `bow.life_sub` を設定（推奨: 100 程度。矢が飛ぶ時間を確保）
   - `script` に着弾時コマンドを記述（`@s` = 被弾エンティティ）
   - `slots` を設定
2. `make run/export` で生成
3. 確認: `generated/passive/effect/{id}.mcfunction` に弓検知コードが生成されていること
4. 確認: `generated/passive/bow/{id}.mcfunction` に Script が出力されていること
5. ゲーム内テスト:
   - `/function maf:generated/passive/give/{id}_slot{N}` で設定書入手
   - 設定書を使用してスロットに装備
   - 弓でモブを撃ち、着弾時にエフェクトが発動することを確認

---

## 8. 関連ファイル一覧

### Generator（Go）
| ファイル | 役割 |
|---------|------|
| `maf-command-generator/app/domain/model/passive/types.go` | Passive / BowConfig 構造体 |
| `maf-command-generator/app/domain/export/passive.go` | BuildPassiveArtifacts, buildBowEffectBody, WritePassiveArtifacts |
| `maf-command-generator/savedata/passive.json` | パッシブマスターデータ |

### Datapack（手書き）
| ファイル | 役割 |
|---------|------|
| `magic/passive/tick.mcfunction` | 毎 tick パッシブ実行（弓 effect もここから呼ばれる） |
| `magic/passive/tag_passive_arrow.mcfunction` | 矢へのデータ埋め込み（マクロ関数） |
| `magic/passive/on_arrow_hit.mcfunction` | 着弾検知（フラグセットのみ） |
| `magic/passive/on_bow_hit.mcfunction` | 着弾処理本体（tick から呼ばれる） |
| `magic/passive/run_bow_effect.mcfunction` | bow スクリプトへのマクロディスパッチ |
| `advancement/arrow_hit.json` | 矢着弾の Advancement トリガー |
| `load.mcfunction` | `mafBowUsed` スコアボード登録 |
| `system/score/afterscore.mcfunction` | `mafBowUsed` リセット |

### Datapack（自動生成）
| ファイル | 役割 |
|---------|------|
| `generated/passive/effect/{id}.mcfunction` | 弓検知 + 矢タグ付けコード |
| `generated/passive/bow/{id}.mcfunction` | 着弾時エフェクト（Script[]） |
| `generated/passive/give/{id}_slot{N}.mcfunction` | 設定書 give |
| `generated/passive/apply/{id}_slot{N}.mcfunction` | スロット書き込み |
