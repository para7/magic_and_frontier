---
name: passive
description: パッシブスキルシステムの設計リファレンス。パッシブのデータモデル・発動条件(always/attack/none)・スロットシステム・生成パイプライン・ランタイム実行フローを網羅する。パッシブの新規追加、発動条件の変更、弓パッシブの調整、装備スロットの操作などで参照すること。passive, パッシブ, スロット, slot, 弓スキル, bow skill, 装備効果, attack, always, none などのキーワードで使う。
---

# パッシブスキルシステム

パッシブは装備や持ち替えで自動発動するスキル。プレイヤーに最大3つのスロットがあり、それぞれに1つのパッシブを装備できる。

## 関連スキル

- **magic-casting**: 詠唱パイプライン。パッシブ「設定書」使用時の装備処理もこれを通る
- **grimoire**: グリモアシステム。詠唱パイプラインを共有
- **ohmydat**: プレイヤー個別ストレージ。パッシブ装備状態の保存先
- **maf-export**: Go ジェネレータのエクスポートパイプライン全体

---

## 1. データモデル（Generator側）

### Passive 構造体

**ファイル:** `maf-command-generator/app/domain/model/passive/types.go`

```go
type Passive struct {
    ID               string   // スラッグID
    Name             string   // 表示名（最大80文字）
    Role             string   // 役割説明（最大200文字）
    Condition        string   // 発動条件: "always" | "attack" | "none"
    Slots            []int    // 装備可能スロット（1〜3, ユニーク, 最低1つ）
    Description      string   // 効果説明（最大400文字）
    Script           []string // 発動時コマンド（1行以上）
    GenerateGrimoire *bool    // true: 設定書を生成しルートテーブルから参照可能 / false: 設定書未生成・ルートテーブル参照不可（必須）
}
```

> 弓スキル（矢への効果）は `Passive` には含まれず、別エンティティ `BowPassive`（`model/bow` パッケージ、`savedata/bow/`）で定義される。詳細は **bow-passive スキル** 参照。

### generate_grimoire フィールド

| 値 | 設定書(give/apply)生成 | ルートテーブル参照 | アイテムへの直接付与 |
|---|---|---|---|
| `true` | される | できる | できる |
| `false` | されない | エラーになる | できる |

- 必須フィールド。省略するとバリデーションエラー
- `false` にすることで「設定書では入手できないが、特定アイテムに付与できるパッシブ」を定義できる

### 発動条件（Condition）の意味

| Condition | 発動タイミング | 処理ファイル |
|-----------|-------------|------------|
| `always` | 毎tick（装備中ずっと） | `passive/tick` → `run_effect` |
| `attack` | 近接攻撃でダメージを与えた時 | `advancement/melee_hit` → `passive/on_melee_hit` → `passive/tick` |
| `none` | 自動発動なし（手動呼び出し専用） | `function maf:generated/passive/effect/{id}` を他システムから実行 |

---

## 2. 生成パイプライン

**ファイル:** `maf-command-generator/app/domain/export/passive.go`

```
savedata/passive/
    ↓ PassiveEntity.Load()
model/passive/Passive[]
    ↓ export.BuildPassiveArtifacts()
    ├── []PassiveEffectFunction   { ID, Body }
    └── []PassiveGrimoireFunction { PassiveID, Slot, FunctionID, GiveBody, ApplyBody, Book }
    ↓ WritePassiveArtifacts()
generated/passive/
    ├── effect/{id}.mcfunction        — 効果スクリプト（`Script[]` を結合）
    ├── give/{id}_slot{N}.mcfunction  — 設定書 give コマンド
    └── apply/{id}_slot{N}.mcfunction — スロット書き込み処理
```

### 生成される成果物

1. **effect/{id}.mcfunction** — `Script[]` をそのまま結合した効果スクリプト
2. **give/{id}_slot{N}.mcfunction** — `GenerateGrimoire=true` のパッシブのみ、Slots ごとに1ファイル
3. **apply/{id}_slot{N}.mcfunction** — `oh_my_dat` にパッシブID/conditionを書き込み + 設定メッセージ表示

> 矢に対して作用する「弓スキル」は `Passive` ではなく **`BowPassive`** エンティティ（`savedata/bow/`、`maf-command-generator/app/domain/export/bow.go`）で実装される。出力される `generated/passive/effect/bow_{id}.mcfunction` / `generated/passive/bow/{id}.mcfunction` は BowPassive の成果物。詳細は **bow-passive スキル**を参照。

---

## 3. ランタイム実行フロー（Datapack側）

### 3a. パッシブ tick（毎tick、全パッシブ共通）

```
maf:magic/tick
  → maf:passive/tick          (全プレイヤー毎tick)
    → oh_my_dat の slot1/2/3 に id があれば → maf:passive/run_effect with slot{N}
    → メインハンドの custom_data から id/condition を oh_my_dat の tmp に書き込み
      → maf:passive/run_effect with tmp
    → mafBowHit >= 1 なら → maf:bow/on_bow_hit（弓パッシブの着弾処理）
```

`run_effect` は `condition` を見て分岐:

| condition | 挙動 |
|-----------|------|
| `always`  | そのまま `generated/passive/effect/$(id)` を呼ぶ |
| `attack`  | `mafMeleeHit >= 1` のときだけ呼ぶ（`advancement/melee_hit` → `passive/on_melee_hit` でセット） |
| `none`    | 何もしない（他 function から手動で呼ぶ用途） |

### 3b. 近接ヒットフラグ（attack condition 用）

```
advancement/melee_hit.json → maf:passive/on_melee_hit
  → advancement revoke + mafMeleeHit += 1
  次の passive/tick → run_effect が attack condition のパッシブを発動
```

### 3c. パッシブ設定書の使用（装備処理）

パッシブ設定書も `minecraft:book` アイテムで、`spell.kind:"passive"` を持つ。  
使用すると通常の詠唱パイプラインを通り、`cast/exec` で `kind:"passive"` を検知:

```
cast/exec → run_passive_apply with {id, slot}
  → $function maf:generated/passive/apply/$(id)_slot$(slot)
    → oh_my_dat の maf.passive.slot{N}.id に passiveID を書き込み
    → "[slotN]に[パッシブ名]を設定しました" メッセージ
```

- 詠唱時間: 固定 200 tick（10秒）
- MP消費: 固定 10 MP
- クールダウン: 0

---

## 4. oh_my_dat ストレージ構造

### パッシブ装備状態

```
oh_my_dat: _[-4]...[-4].maf.passive
├── slot1 { id, condition }   ← スロット1に装備中のパッシブ
├── slot2 { id, condition }   ← スロット2
├── slot3 { id, condition }   ← スロット3
└── tmp  { id, condition }    ← メインハンド用の一時領域（毎tick書き込み→削除）
```

- `slot1`〜`slot3` は永続。設定書を使用すると `maf:generated/passive/apply/{id}_slot{N}` が `id` と `condition` を書き込む
- `tmp` は一時的。`passive/tick` がメインハンドアイテムの `custom_data.maf.passiveId` / `passiveCondition` を毎 tick 書き込み → `run_effect` で使用 → 削除する

### アイテム側のパッシブ情報

```snbt
minecraft:custom_data={
  maf:{
    passiveId: "regeneration",
    passiveSlot: 1,
    passiveCondition: "always",
    passive: {id:"regeneration", slot:1, condition:"always", ...},
    spell: {kind:"passive", id:"regeneration", slot:1, cost:10, cast:200, ...}
  }
}
```

---

## 5. パッシブ追加手順

1. `savedata/passive/entity.json` の `entries` にエントリ追加
2. `condition` を設定（`always`/`attack`/`none`）
3. `slots` を設定（例: `[1, 2]` → スロット1か2に装備可能）
4. `script` に効果コマンドを記述
5. `generate_grimoire` を設定（`true`: 設定書を生成する / `false`: 設定書不要）
6. `none` は tick 自動発動しないため、必要に応じて他functionから `maf:generated/passive/effect/{id}` を呼び出す
7. `make run/export` で生成
8. `generate_grimoire: true` なら `/function maf:generated/passive/give/{id}_slot{N}` でテスト
