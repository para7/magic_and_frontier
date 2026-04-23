---
name: bow-passive
description: 弓パッシブの動作フローリファレンス。BowPassive モデル（bow.json）の4種スクリプト（hit/fired/flying/ground）、矢へのデータ埋め込み、dolphins_grace マーカーによる被弾エンティティ検出、PierceLevel による矢残存、飛翔中・着地後エフェクトなど弓パッシブ固有の仕組みを解説。弓パッシブのバグ修正、新規弓スキル追加、矢の挙動変更、着弾・飛翔・着地エフェクトの調整などで参照すること。bow, 弓, arrow, 矢, 着弾, hit, flying, ground, fired, tag_bow_arrow, on_arrow_hit, mafBowUsed, BowPassive などのキーワードで使う。
---

# 弓パッシブ動作フロー

弓パッシブ（BowPassive）は弓の発射・飛翔・着弾・着地の各フェーズでエフェクトを発動するスキルシステム。`bow.json` で定義され、4種のスクリプトタイプを持つ。

## 関連スキル

- **passive**: パッシブシステム全体。弓パッシブの effect ファイルは passive tick から呼ばれる
- **magic-casting**: 詠唱パイプライン
- **ohmydat**: プレイヤー個別ストレージ
- **maf-export**: Go ジェネレータのエクスポートパイプライン

---

## 1. データモデル（BowPassive）

**ファイル:** `maf-command-generator/app/domain/model/bow/types.go`

```go
type BowPassive struct {
    ID           string   // スラッグID
    Name         string   // 表示名（最大80文字）
    Role         string   // 役割説明（最大200文字）
    Slots        []int    // 装備可能スロット
    LifeSub      *int     // 矢の寿命短縮量（0〜1200 tick）
    ScriptHit    []string // 着弾時スクリプト（@s = 被弾エンティティ）
    ScriptFired  []string // 発射時スクリプト（effect ファイル末尾に追加）
    ScriptFlying []string // 飛翔中スクリプト（矢が飛んでいる間毎 tick）
    ScriptGround []string // 着地後スクリプト（矢が地面に刺さった後毎 tick）
}
```

### 4種のスクリプトタイプ

| スクリプト | 実行タイミング | @s | 生成先 |
|-----------|-------------|-----|--------|
| `ScriptFired` | 弓発射時（effect ファイル末尾） | プレイヤー | `generated/passive/effect/bow_{id}.mcfunction` 内 |
| `ScriptFlying` | 矢飛翔中（毎 tick） | 矢エンティティ | `generated/bow/flying/{id}_flying.mcfunction` |
| `ScriptGround` | 矢着地後（毎 tick） | 矢エンティティ | `generated/bow/ground/{id}_ground.mcfunction` |
| `ScriptHit` | 矢着弾時 | 被弾エンティティ | `generated/passive/bow/{id}.mcfunction` |

### マスターデータ

**ファイル:** `maf-command-generator/savedata/bow.json`

---

## 2. 全体フロー概要

```
[発射フェーズ — 毎 tick]
passive/tick → run_effect → generated/passive/effect/bow_{id}
  → mafBowUsed >= 1 ? (弓を使ったか)
  → YES → 近くの矢を探す → bow/tag_bow_arrow でデータ埋め込み
  → ScriptHit あり → bow/prepare_hit_arrow（dolphins_grace マーカー付与）
  → ScriptFlying あり → tag @s add flying
  → ScriptGround あり → tag @s add ground
  → ScriptFired のコマンドを実行

[飛翔中 — 毎 tick]
magic/tick → bow/tick_flying
  → maf_bow_arrow + flying タグ + 飛行中の矢に対して:
    → bow/run_flying → generated/bow/flying/{id}_flying

[着地後 — 毎 tick]
magic/tick → bow/tick_ground
  → flying タグの除去（地面に刺さった矢）
  → hit タグの除去 + potion_contents 除去（地面に刺さった hit 矢）
  → ground タグの矢に対して:
    → bow/run_ground → generated/bow/ground/{id}_ground

[着弾検知 — advancement トリガー]
advancement/arrow_hit.json (player_hurt_entity)
  → passive/on_arrow_hit
    → advancement revoke（再発火可能に）
    → mafBowHit += 1（フラグのみ。処理は tick に委譲）

[着弾処理 — 次 tick]
passive/tick → mafBowHit >= 1 → bow/on_bow_hit
  → maf_bow_arrow タグ + shooterPlayerID 一致の矢を検索
  → bow/resolve_hit_arrow:
    → dolphins_grace(amp:80) を持つエンティティ（被弾者）に対して:
       → bow/run_bow_effect → generated/passive/bow/{id}
    → dolphins_grace マーカー除去
    → 矢を kill
```

---

## 3. 発射フェーズ（矢のタグ付け）

### 3a. effect ファイル（自動生成）

**生成元:** `maf-command-generator/app/domain/export/bow.go` の `buildBowPassiveEffectBody()`

**生成先:** `datapacks/.../generated/passive/effect/bow_{id}.mcfunction`

```mcfunction
execute unless score @s mafBowUsed matches 1.. run return 0
execute store result storage maf:tmp bow_player_id int 1 run scoreboard players get @s mafPlayerID
execute as @e[type=arrow,distance=..2,nbt=!{inGround:1b},sort=nearest,limit=1] run function maf:bow/tag_bow_arrow {bow_id:"{id}",life:{lifeValue}}
execute as @e[type=arrow,distance=..2,tag=maf_bow_arrow,sort=nearest,limit=1] run function maf:bow/prepare_hit_arrow  # ScriptHit ありの場合のみ
execute as @e[type=arrow,distance=..2,tag=maf_bow_arrow,sort=nearest,limit=1] run tag @s add flying                    # ScriptFlying ありの場合のみ
execute as @e[type=arrow,distance=..2,tag=maf_bow_arrow,sort=nearest,limit=1] run tag @s add ground                    # ScriptGround ありの場合のみ
{ScriptFired のコマンド群}
```

| 行 | 処理 | 条件 |
|----|------|------|
| mafBowUsed チェック | 弓使用がなければ即 return | 常に |
| bow_player_id 退避 | shooterPlayerID 書き込み用 | 常に |
| tag_bow_arrow | 矢にタグ + データ埋め込み | 常に |
| prepare_hit_arrow | hit タグ + dolphins_grace マーカー | ScriptHit 非空時 |
| tag flying | 飛翔エフェクト対象フラグ | ScriptFlying 非空時 |
| tag ground | 着地エフェクト対象フラグ | ScriptGround 非空時 |
| ScriptFired | 発射時の追加コマンド | ScriptFired 非空時 |

### 3b. mafBowUsed スコアボード

| ファイル | 処理 |
|---------|------|
| `load.mcfunction` | `scoreboard objectives add mafBowUsed minecraft.used:minecraft.bow` — バニラ統計から自動カウント |
| `system/score/afterscore.mcfunction` | `scoreboard players set @a mafBowUsed 0` — 毎 tick 末にリセット |

弓を使った tick だけ `mafBowUsed >= 1` になる。次 tick にはリセットされるため、1回の射撃で1回だけ矢タグ付けが走る。

### 3c. bow/tag_bow_arrow.mcfunction（矢へのデータ埋め込み）

**ファイル:** `datapacks/.../bow/tag_bow_arrow.mcfunction`

```mcfunction
# @s = arrow entity, args: bow_id, life
$data merge entity @s {Tags:["maf_bow_arrow"],PierceLevel:2b,item:{components:{"minecraft:custom_data":{maf:{bowId:"$(bow_id)"}}}}}
$data modify entity @s life set value $(life)s
data modify entity @s item.components."minecraft:custom_data".maf.shooterPlayerID set from storage maf:tmp bow_player_id
```

矢エンティティに書き込む情報:

| データ | 値 | 目的 |
|--------|-----|------|
| `Tags: ["maf_bow_arrow"]` | 固定 | 弓パッシブ矢を `@e` セレクタで検索 |
| `PierceLevel: 2b` | 固定 | 命中しても矢が消えないようにする（着弾処理後に手動 kill） |
| `custom_data.maf.bowId` | 弓パッシブID | 着弾時・飛翔時・着地時にどのスキルか識別 |
| `life` | `1200 - LifeSub` | 矢の生存 tick |
| `custom_data.maf.shooterPlayerID` | プレイヤーID | 着弾時にシューターの矢だけを処理するため |

### 3d. bow/prepare_hit_arrow.mcfunction（hit マーカー付与）

**ファイル:** `datapacks/.../bow/prepare_hit_arrow.mcfunction`

```mcfunction
tag @s add hit
data merge entity @s {item:{components:{"minecraft:potion_contents":{custom_effects:[{id:"minecraft:dolphins_grace",duration:200,amplifier:80}]}}}}
```

- `hit` タグ: 着弾処理対象であることを示す
- dolphins_grace (amp:80): 矢が当たったエンティティに付与される被弾マーカー

### 3e. life 値の計算

```go
lifeSub := 1200  // デフォルト: 即消滅
if entry.LifeSub != nil {
    lifeSub = *entry.LifeSub
}
lifeValue := 1200 - lifeSub
```

| LifeSub | life 値 | 矢の飛行時間 |
|---------|---------|-------------|
| 未設定 (nil) | 0 | 即消滅 |
| 0 | 1200 | 60秒（バニラ標準） |
| 100 | 1100 | 55秒 |
| 1200 | 0 | 即消滅 |

> **注意:** LifeSub 未設定時のデフォルトは 1200（life=0、即消滅）。実用的な弓パッシブでは必ず `LifeSub` を明示的に指定すること。

---

## 4. 飛翔フェーズ（flying）

### 4a. bow/tick_flying.mcfunction

**ファイル:** `datapacks/.../bow/tick_flying.mcfunction`

```mcfunction
execute as @e[type=arrow,tag=maf_bow_arrow,tag=flying,nbt=!{inGround:1b}] at @s run function maf:bow/run_flying with entity @s item.components."minecraft:custom_data".maf
```

`magic/tick.mcfunction` から直接呼ばれる（プレイヤーコンテキスト不要）。飛行中（`inGround:0b`）かつ `flying` タグ付きの弓矢に対して、矢の custom_data から `bowId` をマクロ引数として `run_flying` を呼ぶ。

### 4b. bow/run_flying.mcfunction

```mcfunction
$function maf:generated/bow/flying/$(bowId)_flying
```

`@s` は矢エンティティ。ScriptFlying のコマンドは矢の位置で実行される。

---

## 5. 着地フェーズ（ground）

### 5a. bow/tick_ground.mcfunction

**ファイル:** `datapacks/.../bow/tick_ground.mcfunction`

```mcfunction
# flying タグの矢が地面に刺さったら flying を除去
execute as @e[type=arrow,tag=maf_bow_arrow,tag=flying,nbt={inGround:1b}] run tag @s remove flying
# hit タグの矢が地面に刺さったら potion_contents 除去 + hit タグ除去
execute as @e[type=arrow,tag=maf_bow_arrow,tag=hit,nbt={inGround:1b}] run data remove entity @s item.components."minecraft:potion_contents"
execute as @e[type=arrow,tag=maf_bow_arrow,tag=hit,nbt={inGround:1b}] run tag @s remove hit
# ground タグの矢が地面に刺さったらエフェクト実行
execute as @e[type=arrow,tag=maf_bow_arrow,tag=ground,nbt={inGround:1b}] at @s run function maf:bow/run_ground with entity @s item.components."minecraft:custom_data".maf
```

`magic/tick.mcfunction` から直接呼ばれる。地面に刺さった矢に対する処理:
1. `flying` タグ除去（飛翔エフェクト停止）
2. `hit` タグ + `potion_contents` 除去（地面に刺さった hit 矢の dolphins_grace を無効化）
3. `ground` タグの矢に着地エフェクト実行

### 5b. bow/run_ground.mcfunction

```mcfunction
$function maf:generated/bow/ground/$(bowId)_ground
```

`@s` は矢エンティティ。ScriptGround のコマンドは矢の位置で実行される。

---

## 6. 着弾フェーズ（hit）

着弾処理は **2段階**に分離されている:
1. advancement が即座にフラグを立てる（`passive/on_arrow_hit`）
2. 次の tick ループで実際の処理を行う（`bow/on_bow_hit`）

### 6a. Advancement トリガー

**ファイル:** `datapacks/.../advancement/arrow_hit.json`

```json
{
  "rewards": {
    "function": "maf:passive/on_arrow_hit"
  }
}
```

- `player_hurt_entity` + `direct_entity: arrow` — プレイヤーが矢でエンティティを傷つけたとき発火
- `@s` は **攻撃側プレイヤー**

### 6b. passive/on_arrow_hit.mcfunction

**ファイル:** `datapacks/.../passive/on_arrow_hit.mcfunction`

```mcfunction
advancement revoke @s only maf:arrow_hit
scoreboard players add @s mafBowHit 1
```

フラグセットのみ。実処理は tick に委譲。

### 6c. bow/on_bow_hit.mcfunction（着弾処理本体）

**ファイル:** `datapacks/.../bow/on_bow_hit.mcfunction`

```mcfunction
scoreboard players set @s mafBowHit 0
execute store result storage maf:tmp bow_player_id int 1 run scoreboard players get @s mafPlayerID
function maf:bow/process_hit_arrows with storage maf:tmp
data remove storage maf:tmp bow_player_id
```

`passive/tick` から `mafBowHit >= 1` のとき呼ばれる。プレイヤーの `mafPlayerID` を使って、自分が撃った矢だけを処理する。

### 6d. bow/process_hit_arrows.mcfunction

```mcfunction
$execute as @e[type=arrow,tag=maf_bow_arrow,nbt={item:{components:{"minecraft:custom_data":{maf:{shooterPlayerID:$(bow_player_id)}}}}}] at @s run function maf:bow/resolve_hit_arrow
```

`shooterPlayerID` が一致する `maf_bow_arrow` タグ付き矢を検索し、各矢に対して `resolve_hit_arrow` を実行。

### 6e. bow/resolve_hit_arrow.mcfunction

**ファイル:** `datapacks/.../bow/resolve_hit_arrow.mcfunction`

```mcfunction
execute unless entity @e[nbt={active_effects:[{id:"minecraft:dolphins_grace",amplifier:80b}]},sort=nearest,limit=1,distance=..5] run return 0
data modify storage maf:tmp bow_passive set from entity @s item.components."minecraft:custom_data".maf
execute if entity @s[tag=hit] as @e[nbt={active_effects:[{id:"minecraft:dolphins_grace",amplifier:80b}]},sort=nearest,limit=1,distance=..5] at @s run function maf:bow/run_bow_effect with storage maf:tmp bow_passive
effect clear @e[nbt={active_effects:[{id:"minecraft:dolphins_grace",amplifier:80b}]},sort=nearest,limit=1,distance=..5] minecraft:dolphins_grace
kill @s
data remove storage maf:tmp bow_passive
```

| ステップ | 処理 | なぜ必要か |
|---------|------|----------|
| dolphins_grace チェック | 近くに被弾マーカー持ちがいなければ return | hit 以外の矢（flying/ground のみ）でも呼ばれうるため |
| データ退避 | 矢の `custom_data.maf` を `maf:tmp bow_passive` にコピー | kill 後はデータ取得不可 |
| hit チェック + エフェクト実行 | `hit` タグ付き矢の場合のみ被弾エンティティでエフェクト実行 | ScriptHit がある矢だけ |
| マーカー除去 | dolphins_grace (amp:80) を除去 | 副作用防止 |
| 矢 kill | 矢エンティティを除去 | PierceLevel で残存した矢の後始末 |
| クリア | 一時ストレージ削除 | ゴミデータ残留防止 |

### 6f. 被弾エンティティの検出方式（dolphins_grace マーカー方式）

矢の `potion_contents` に仕込んだ dolphins_grace (amplifier:80) が命中時に被弾エンティティへ付与される。このエフェクトを持つエンティティ＝被弾エンティティとして正確に特定できる。

| 特性 | 内容 |
|------|------|
| 特定精度 | 高い。amplifier:80 は通常付与されない特殊値 |
| 複数命中 | PierceLevel:2b により貫通。全被弾エンティティにエフェクト発動 |
| クリア必須 | 泳ぎ速度上昇の副作用を防ぐため必ず `effect clear` |

### 6g. bow/run_bow_effect.mcfunction

```mcfunction
$function maf:generated/passive/bow/$(bowId)
```

`maf:tmp bow_passive` ストレージから `bowId` を受け取り、マクロで `generated/passive/bow/{id}.mcfunction` にディスパッチ。`@s` は被弾エンティティ。

---

## 7. 生成パイプライン（Go 側）

### 7a. 成果物の分離

**ファイル:** `maf-command-generator/app/domain/export/bow.go`

```
BuildBowArtifacts(master)
  → []BowEffectFunction   — effect ファイル（弓検知 + 矢タグ付け + ScriptFired）
  → []BowHitFunction      — hit ファイル（ScriptHit）
  → []BowFlyingFunction   — flying ファイル（ScriptFlying）
  → []BowGroundFunction   — ground ファイル（ScriptGround）
```

| 成果物 | パス | 内容 | ID 形式 |
|--------|------|------|---------|
| effect | `generated/passive/effect/bow_{id}.mcfunction` | 弓検知 + 矢タグ付け + ScriptFired | `bow_` プレフィックス付き |
| hit | `generated/passive/bow/{id}.mcfunction` | 着弾時エフェクト | そのまま |
| flying | `generated/bow/flying/{id}_flying.mcfunction` | 飛翔中エフェクト | `_flying` サフィックス |
| ground | `generated/bow/ground/{id}_ground.mcfunction` | 着地エフェクト | `_ground` サフィックス |

**注意:** ScriptHit / ScriptFlying / ScriptGround が空の場合、対応するファイルは生成されない（既存ファイルも削除される）。

### 7b. ファイル書き出し

**ファイル:** `maf-command-generator/app/domain/export/bow.go` の `WriteBowArtifacts()`

effect と hit は passive と同じディレクトリ（`passiveEffectDir`, `passiveBowDir`）に出力される。flying と ground は `generated/bow/` 配下の専用ディレクトリに出力される。

---

## 8. magic/tick.mcfunction からの呼び出し

```mcfunction
execute as @a at @s run function maf:passive/tick      # パッシブ tick（effect 実行 + mafBowHit チェック）
function maf:bow/tick_flying                            # 飛翔中エフェクト（全矢対象）
function maf:bow/tick_ground                            # 着地エフェクト（全矢対象）
```

`tick_flying` と `tick_ground` はプレイヤーコンテキスト不要なのでグローバルに呼ばれる。

---

## 9. 設計上の注意点

### PierceLevel の必要性

通常の矢は命中で消滅する。矢が消えると `custom_data` が失われるため、`PierceLevel:2b` で貫通属性を付けて矢を残存させる。着弾処理後に `resolve_hit_arrow` で手動 kill する。

### dolphins_grace マーカー方式の注意点

- `hit` タグのない矢（flying/ground のみ）には `prepare_hit_arrow` が呼ばれないため、dolphins_grace は付与されない
- `tick_ground` で地面に刺さった `hit` 矢の `potion_contents` を除去するため、地面への着弾で誤検出は起きない
- `effect clear` を忘れるとエンティティに泳ぎ速度上昇が残り続ける

### shooterPlayerID によるフィルタリング

`bow/process_hit_arrows` で `shooterPlayerID` が一致する矢だけを処理する。複数プレイヤーが弓を使った場合でも、各プレイヤーは自分の矢だけを処理する。

### LifeSub デフォルト値の罠

`LifeSub` が nil（未設定）の場合、デフォルトの lifeSub は **1200** になり、`life = 0` で矢が即消滅する。新規弓パッシブには必ず `LifeSub` を設定すること。

---

## 10. 弓パッシブ追加手順

1. `savedata/bow.json` にエントリ追加
   - `life_sub` を設定（推奨: 100 程度）
   - 必要なスクリプトを記述:
     - `script_hit`: 着弾時コマンド（`@s` = 被弾エンティティ）
     - `script_fired`: 発射時追加コマンド（`@s` = プレイヤー）
     - `script_flying`: 飛翔中コマンド（`@s` = 矢、毎 tick）
     - `script_ground`: 着地後コマンド（`@s` = 矢、毎 tick）
2. `make run/export` で生成
3. 確認: 各スクリプトに対応するファイルが生成されていること
4. ゲーム内テスト: 弓パッシブ付き武器で各フェーズの動作を確認

---

## 11. 関連ファイル一覧

### Generator（Go）
| ファイル | 役割 |
|---------|------|
| `maf-command-generator/app/domain/model/bow/types.go` | BowPassive 構造体 |
| `maf-command-generator/app/domain/export/bow.go` | BuildBowArtifacts, buildBowPassiveEffectBody, WriteBowArtifacts |
| `maf-command-generator/savedata/bow.json` | 弓パッシブマスターデータ |

### Datapack（手書き）
| ファイル | 役割 |
|---------|------|
| `bow/tag_bow_arrow.mcfunction` | 矢へのデータ埋め込み（マクロ関数） |
| `bow/prepare_hit_arrow.mcfunction` | hit マーカー（hit タグ + dolphins_grace）付与 |
| `bow/on_bow_hit.mcfunction` | 着弾処理エントリ（フラグリセット + process_hit_arrows 呼び出し） |
| `bow/process_hit_arrows.mcfunction` | shooterPlayerID で矢をフィルタして resolve |
| `bow/resolve_hit_arrow.mcfunction` | 個別矢の着弾解決（エフェクト実行 + クリア + kill） |
| `bow/run_bow_effect.mcfunction` | hit スクリプトへのマクロディスパッチ |
| `bow/tick_flying.mcfunction` | 飛翔中矢の毎 tick 処理 |
| `bow/run_flying.mcfunction` | flying スクリプトへのマクロディスパッチ |
| `bow/tick_ground.mcfunction` | 着地矢の毎 tick 処理 |
| `bow/run_ground.mcfunction` | ground スクリプトへのマクロディスパッチ |
| `passive/on_arrow_hit.mcfunction` | 着弾検知（advancement コールバック、フラグセットのみ） |
| `advancement/arrow_hit.json` | 矢着弾の Advancement トリガー |

### Datapack（自動生成）
| ファイル | 役割 |
|---------|------|
| `generated/passive/effect/bow_{id}.mcfunction` | 弓検知 + 矢タグ付け + ScriptFired |
| `generated/passive/bow/{id}.mcfunction` | 着弾時エフェクト（ScriptHit） |
| `generated/bow/flying/{id}_flying.mcfunction` | 飛翔中エフェクト（ScriptFlying） |
| `generated/bow/ground/{id}_ground.mcfunction` | 着地エフェクト（ScriptGround） |
