---
name: sight
description: 視線判定（ライン・オブ・サイト）システムの設計リファレンス。プレイヤーの視線方向にいるモブを検出してタグ付けする共通ユーティリティ。グリモアのエフェクト Script[] にこの仕組みを使いたい場合、パッシブで視線方向の敵を対象にしたい場合、または「視線上の敵」「見ている方向」「ターゲット指定」系の処理を書くときは必ずこのスキルを参照すること。
---

# 視線判定システム（Line of Sight）

プレイヤーの目線方向に0.5ブロック刻みでレイキャストを行い、視線上の敵モブにタグを付与する共通ユーティリティ。

## 関連スキル

- **grimoire**: グリモアのエフェクト Script[] でよく使う
- **passive**: パッシブスキル効果でも利用可能

---

## 1. 関数一覧

| 関数 | ファイル | 特徴 |
|------|---------|------|
| `maf:common/sight/eyes_tagged` | `common/sight/eyes_tagged.mcfunction` | 通常視線。水・溶岩・不透過ブロックで遮断 |
| `maf:common/sight/eyes_tagged_water` | `common/sight/eyes_tagged_water.mcfunction` | 水中視線。水を透過し溶岩・不透過ブロックのみ遮断 |
| `maf:common/sight/eyes_tagged_through_blocks` | `common/sight/eyes_tagged_through_blocks.mcfunction` | ブロック無視視線。すべてのブロックを透過してモブを検出 |

すべて **呼び出し元コンテキストがプレイヤー（`@s`）** であること前提。

---

## 2. 付与されるタグ

| タグ | 対象 | 内容 |
|------|------|------|
| `maf_sight_target` | 最大1体 | 視線上で最初に検出された最近傍のモブ |
| `maf_sight_candidate` | 0体以上 | 視線上の全ステップで検出されたモブ（複数可） |

- `maf_sight_target` は「主要ターゲット1体」（ダメージ・移動など単体対象の処理に使う）
- `maf_sight_candidate` は「視線ライン上の全員」（グローなどAOE風の処理に使う）

---

## 3. レイキャストの仕組み

0.5ブロック刻みで最大30ブロック先まで走査。各ステップで:

1. 溶岩ブロックを検出 → `return 0`（即終了）
2. 水ブロックを検出 → `return 0`（`eyes_tagged_water` は水を透過するので別途プレディケートで判定）
3. `#maf:sight_passable` に含まれないブロックを検出 → `return 0`
4. 半径2.0以内に `#maf:enemymob` が存在 → `maf_sight_candidate` タグを付与
5. まだ `maf_sight_target` が存在しなければ → 最近傍の1体に `maf_sight_target` タグを付与

`eyes_tagged_through_blocks` はステップ 1〜3 のブロックチェックをすべて省略し、4〜5 のモブ検出のみ行う。壁越しや地形貫通の魔法に使う。

### ブロック透過設定

| タグ/プレディケート | ファイル | 内容 |
|--------------------|---------|------|
| `#maf:sight_passable` | `tags/block/sight_passable.json` | 通常視線が透過するブロック（`#maf:air` + `#minecraft:replaceable`） |
| `#maf:sight_passable_water` | `tags/block/sight_passable_water.json` | 水中視線が透過するブロック（`#maf:sight_passable` + `#maf:water`） |
| `maf:predicate/sight_water_fluid` | `predicate/sight_water_fluid.json` | 水フルイドを判定するプレディケート |

透過対象を増やすには `sight_passable.json` に `"replace": false` でブロックを追加する。

---

## 4. 使用パターン

視線判定を使うエフェクト関数の標準的な書き方:

```mcfunction
# 1. 視線判定を実行（@s = プレイヤー）
#    タグのクリアは関数内部の先頭で行われる
function maf:common/sight/eyes_tagged

# 2. タグを使って効果を適用
#   単体対象（最近傍1体）
execute as @e[type=#maf:enemymob,tag=maf_sight_target,distance=..8,sort=nearest,limit=1] run damage @s 20 minecraft:player_attack
#   全候補（視線ライン上の全員）
effect give @e[type=#maf:enemymob,tag=maf_sight_candidate,distance=..20] minecraft:glowing 1 0 true
```

タグのクリアは各視線関数の先頭で自動的に行われるため、呼び出し元での前後クリアは不要。

### 水中で使う場合

```mcfunction
function maf:common/sight/eyes_tagged_water
# ... 以降同じ
```

### ブロックを無視する場合

```mcfunction
function maf:common/sight/eyes_tagged_through_blocks
# ... 以降同じ
```

---

## 5. 実装例

### tomahawk01（単体ダメージ、射程8）

```mcfunction
function maf:common/sight/eyes_tagged
execute as @e[type=#maf:enemymob,tag=maf_sight_target,distance=..8,sort=nearest,limit=1] run damage @s 20 minecraft:player_attack
execute as @e[type=#maf:enemymob,tag=maf_sight_target,distance=..8,sort=nearest,limit=1] at @s run particle minecraft:crit ~ ~0.9 ~ 0.3 0.5 0.3 0.01 20 force
```

### sight_levitation01（単体浮遊 + 全候補グロー、水透過）

```mcfunction
function maf:common/sight/eyes_tagged_water
execute as @e[type=#maf:enemymob,tag=maf_sight_target,distance=..20,sort=nearest,limit=1] run effect give @s minecraft:levitation 1 1
effect give @e[type=#maf:enemymob,tag=maf_sight_candidate,distance=..20] minecraft:glowing 1 0 true
```

---

## 6. グリモアの Script[] に組み込む場合

グリモアのエフェクトは生成時に `generated/grimoire/effect/{id}.mcfunction` に書き出されるため、`Script[]` に直接 `function maf:common/sight/eyes_tagged` の呼び出しを含める。

```json
{
  "id": "my_spell01",
  "script": [
    "function maf:common/sight/eyes_tagged",
    "execute as @e[type=#maf:enemymob,tag=maf_sight_target,distance=..15,sort=nearest,limit=1] run effect give @s minecraft:wither 3 1"
  ]
}
```

### 注意点

- タグのクリアは各視線関数の **先頭** で行われる。呼び出し元での手動クリアは不要。
- `distance=..N` の射程フィルタは視線関数の最大30ブロックより短く設定してよい（魔法ごとの射程で制御）。
- `sort=nearest,limit=1` を忘れると複数のターゲットに効果が発動する場合がある。
