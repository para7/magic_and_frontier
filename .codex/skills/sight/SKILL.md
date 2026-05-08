---
name: sight
description: 視線判定（ライン・オブ・サイト）システムの設計リファレンス。プレイヤーの視線方向にいるモブを検出してタグ付けする共通ユーティリティ。アクティブのエフェクト Script[] にこの仕組みを使いたい場合、パッシブで視線方向の敵を対象にしたい場合、または「視線上の敵」「見ている方向」「ターゲット指定」系の処理を書くときは必ずこのスキルを参照すること。
---

# 視線判定システム（Line of Sight）

プレイヤーの視線方向にいる敵モブにタグを付与する共通ユーティリティ。呼び出し元は必ずプレイヤー（`@s`）であること。

## 関数一覧

| 関数 | 用途 |
|------|------|
| `maf:common/sight/eyes_tagged` | 通常視線。水・溶岩・不透過ブロックで遮断 |
| `maf:common/sight/eyes_tagged_water` | 水中視線。水を透過し溶岩・不透過ブロックのみ遮断 |
| `maf:common/sight/eyes_tagged_through_blocks` | ブロック無視視線。壁越しや地形貫通の魔法に使う |

## 付与されるタグ

| タグ | 対象 | 用途 |
|------|------|------|
| `maf_target` | 最大1体 | 視線上の最近傍モブ。単体ダメージ・移動など単体処理に使う |
| `maf_sight_candidate` | 0体以上 | 視線ライン上の全モブ。範囲処理に使う |

タグのクリアは各視線関数の先頭で自動的に行われるため、呼び出し元での手動クリアは不要。

## 使用例

### 単体ダメージ（射程8）

```mcfunction
function maf:common/sight/eyes_tagged
execute as @e[type=#maf:enemymob,tag=maf_target,distance=..8,sort=nearest,limit=1] run damage @s 20 minecraft:player_attack
```

### 単体浮遊 + 視線上全員グロー（水透過、射程20）

```mcfunction
function maf:common/sight/eyes_tagged_water
execute as @e[type=#maf:enemymob,tag=maf_target,distance=..20,sort=nearest,limit=1] run effect give @s minecraft:levitation 1 1
effect give @e[type=#maf:enemymob,tag=maf_sight_candidate,distance=..20] minecraft:glowing 1 0 true
```

### アクティブの Script[] に組み込む場合

```json
{
  "id": "my_spell01",
  "script": [
    "function maf:common/sight/eyes_tagged",
    "execute as @e[type=#maf:enemymob,tag=maf_target,distance=..15,sort=nearest,limit=1] run effect give @s minecraft:wither 3 1"
  ]
}
```

視線処理は固定で30ブロックにタグを付与するので、 `distance=..N` の射程フィルタは呼び出し側で制御する。
