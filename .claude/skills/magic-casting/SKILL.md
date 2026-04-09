---
name: magic-casting
description: 魔法詠唱パイプラインの設計リファレンス。キャスト・クールダウン・MP消費・MPバー表示の共通基盤で、グリモアとパッシブの両方が通る。スコアボード、oh_my_dat ストレージ、ソウルシステムを含む。詠唱関連のバグ修正、MP/ソウルの挙動変更、スコアボード調査、MPバー表示の修正などで参照すること。cast, 詠唱, MP, マナ, クールダウン, cooldown, ソウル, soul, スコアボード, scoreboard, MPバー などのキーワードで使う。
---

# 魔法詠唱パイプライン

グリモアとパッシブ（設定書使用時）が共有する詠唱→発動→クールダウンの仕組み。

## 関連スキル

- **grimoire**: グリモア固有のデータモデル・生成・NBT構造
- **passive**: パッシブ固有のデータモデル・発動条件
- **ohmydat**: `oh_my_dat` ストレージのアクセス方法

---

## 1. 全体フロー

```
プレイヤーが本を右クリック
    ↓
advancement/use_grimoire.json (using_item トリガー)
    ↓
maf:magic/use_grimoire (状態チェック → spell データを oh_my_dat にコピー)
    ↓
maf:magic/exec/set_magic (ストレージ → スコアボードにロード、MP検証)
    ↓
毎tick: maf:magic/cast/tick (カウントダウン、移動キャンセル判定)
    ↓ mafCastTime が 0 になったら
maf:magic/cast/exec (MP消費 → kind で分岐ディスパッチ)
    ├─ kind:"grimoire" → run_grimoire_effect → generated/grimoire/effect/{id}
    └─ kind:"passive"  → run_passive_apply   → generated/passive/apply/{id}_slot{N}
    ↓
casting データ削除、クールダウン開始
```

---

## 2. スコアボード一覧

### 魔法システム用（`magic/load.mcfunction` で登録）

| スコアボード | 型 | 用途 | 値の意味 |
|------------|-----|------|---------|
| `mafMP` | dummy | 現在MP | 0〜mafMaxMP |
| `mafMaxMP` | dummy | 最大MP | = mafSoul |
| `mafCastTime` | dummy | 詠唱カウントダウン | **-1**=非詠唱, **0**=発動タイミング, **1+**=詠唱中 |
| `mafCastTimeMax` | dummy | 詠唱時間の初期値 | 進捗バー表示用 |
| `mafCoolTime` | dummy | クールダウン残り | 0以下で使用可能 |
| `mafCastCost` | dummy | 現在のスペルMP消費 | set_magic でロード |
| `mafMPTick` | dummy | MP回復タイマー | 600到達で+1MP、リセット |

### mafCastTime の値が特に重要

- **-1**: 初期状態・非詠唱。`use_grimoire` は `-1以下` のとき受付
- **0**: このtickで `cast/exec` が呼ばれて発動する
- **1以上**: 詠唱中。毎tickデクリメント
- set_magic の冒頭で一度 `-1` にセットしてからロードする（ロード失敗時のバグ対策）

---

## 3. 各ファイルの詳細

### set_magic.mcfunction（パラメータロード）

**ファイル:** `datapacks/.../magic/exec/set_magic.mcfunction`

```
1. mafCastTime = -1 にセット（安全弁）
2. oh_my_dat ストレージからスコアボードにロード:
   - mafCastCost  ← casting.cost
   - mafCastTime  ← casting.cast
   - mafCastTimeMax ← casting.cast
   - mafCoolTime  ← casting.cooltime
3. MP不足チェック: mafCastCost > mafMP なら
   → mafCastTime = -1 に戻す
   → "MPが足りません！" メッセージ
   → casting データ削除
   → return fail
4. 詠唱名を MPバーストレージに保存（p7:mpbar bar{N}.title）
5. return 1
```

### cast/tick.mcfunction（毎tick処理）

**ファイル:** `datapacks/.../magic/cast/tick.mcfunction`

```
1. mafMPTick を -40 にリセット（詠唱中はMP回復停止）
2. 詠唱エフェクト: mafCastTime >= 40 でエンチャントパーティクル
3. 移動キャンセル: mafCastTime >= 11 かつ mafMoved >= 1 → cancel
4. mafCastTime == 0 で exec を呼び出し
5. mafCastTime をデクリメント
6. p7_setSkEnable = -1（スキル無効化）
```

### cast/exec.mcfunction（発動）

**ファイル:** `datapacks/.../magic/cast/exec.mcfunction`

```
1. mafMP -= mafCastCost（MP消費）
2. oh_my_dat から kind を読み取り:
   - "grimoire" → run_grimoire_effect with storage（マクロディスパッチ）
   - "passive"  → run_passive_apply with storage（マクロディスパッチ）
3. casting データ削除
```

### cast/cancel.mcfunction（詠唱中断）

```
1. "詠唱が中断されました" メッセージ
2. mafCastCost=0, mafCastTime=-1, mafCastTimeMax=0
3. casting データ削除
```

---

## 4. oh_my_dat ストレージ構造

詠唱中のスペルデータは以下に保存される:

```
storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting
├── kind: "grimoire" | "passive"
├── id: "healing01"
├── cost: 13
├── cast: 40
├── cooltime: 20
├── title: "ヒーリング"
├── description: "周囲に即時回復+リジェネ"
└── slot: 1  (passive のみ)
```

このデータは `spell` フラグメントとしてアイテムの `custom_data.maf.spell` に埋め込まれ、使用時にそのままコピーされる。

### spell フラグメントの統一性

グリモアもパッシブも同じ `spell` フラグメント形式を使う。違いは:
- グリモア: `kind:"grimoire"`, `slot` なし
- パッシブ: `kind:"passive"`, `slot` あり（1〜3）

この統一により、`set_magic` 以降のパイプラインは kind を見るだけで両方に対応できる。

---

## 5. MP / ソウルシステム

### MP回復（mp_manage.mcfunction）

**ファイル:** `datapacks/.../magic/mp/mp_manage.mcfunction`

```
毎tick: mafMPTick += 10
mafMPTick >= 600 → mafMP += 1, mafMPTick = 0
mafMaxMP = mafSoul（毎tick同期）
mafMP = min(mafMP, mafMaxMP)（キャップ）
```

- **回復速度**: 600 tick = 30秒 で 1MP 回復
- `mafMPTick += 10` なので、内部的には tick あたり10カウント進む
- 詠唱中は `cast/tick` で `mafMPTick = -40` にリセットされるため回復停止

### ソウルシステム（soul/tick.mcfunction）

**ファイル:** `datapacks/.../soul/tick.mcfunction`

```
毎tick: mafSoulTick += 10
mafSoulTick >= 1200 → mafSoul += 1, mafSoulTick = 0
mafSoul は最大100にキャップ
死亡時（mafSoulReset >= 1）→ mafSoul = 0, mafMPTick = 0
```

- **回復速度**: 1200 tick = 60秒 で 1ソウル回復
- ソウル = 最大MP。ソウルが増えると使えるスペルが増える
- 死亡でソウルリセット（ペナルティ）

---

## 6. MPバー表示

**ストレージ:** `p7:mpbar`

プレイヤーIDごとに `bar1`〜`bar20` のスロットがある。

```
storage p7:mpbar bar{N}
├── title: "ヒーリング"  (詠唱中のスペル名)
└── (その他表示データ)
```

- `set_magic` で詠唱開始時にスペル名を書き込む
- `mafPlayerID`（1〜20）で対応するスロットを決定
- 全プレイヤー分をハードコードで書き出している（マクロ展開不可のため）

---

## 7. メインティックループ

**ファイル:** `datapacks/.../magic/tick.mcfunction`

```mcfunction
# クールダウン減算
scoreboard players remove @a[scores={mafCoolTime=1..}] mafCoolTime 1

# 詠唱処理（mafCastTime >= 0 のプレイヤーのみ）
execute as @a at @s if score @s mafCastTime matches 0.. run function maf:magic/cast/tick

# パッシブ効果（全プレイヤー）
execute as @a at @s run function maf:magic/passive/tick

# MP回復
execute as @a at @s run function maf:magic/mp/mp_manage

# MPバー更新
function maf:magic/mp/mpbar
```

処理順序が重要: クールダウン → 詠唱 → パッシブ → MP → 表示
