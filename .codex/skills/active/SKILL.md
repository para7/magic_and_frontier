---
name: active
description: アクティブ（アクティブスキル）システムの設計リファレンス。アクティブのデータモデル・NBT構造・生成パイプライン・ランタイム実行フローを網羅する。アクティブの新規追加、エフェクト修正、アイテムへのバインド、スペルのデバッグなど、アクティブに関わる作業全般で参照すること。active, アクティブスキル, スペル, 魔法書, 呪文, spell, cast などのキーワードが出たら使う。
---

# アクティブシステム

アクティブはプレイヤーが使う「アクティブスキル」アイテム。本（`minecraft:book`）として配布され、右クリックで詠唱→効果発動する。

## 関連スキル

- **magic-casting**: 詠唱パイプライン（キャスト・クールダウン・MP消費の共通基盤）。アクティブとパッシブの両方が使う
- **passive**: パッシブスキルシステム。詠唱パイプラインを共有するがトリガー方式が異なる
- **ohmydat**: プレイヤー個別ストレージ。詠唱データの一時保存先
- **maf-export**: Go ジェネレータのエクスポートパイプライン全体

---

## 1. データモデル（Generator側）

### Active 構造体

**ファイル:** `maf-command-generator/app/domain/model/active/types.go`

```go
type Active struct {
    ID          string   // スラッグID（小文字・ハイフン・アンダースコアのみ）
    CastTime    int      // 詠唱時間（tick単位, 0〜12000）
    CoolTime    int      // クールダウン（tick単位, 0〜12000）
    MPCost      int      // MP消費量（0〜1,000,000）
    Script      []string // 発動時に実行するmcfunctionコマンド群（1行以上必須）
    Title       string   // 表示名（必須）
    Description string   // 説明文（任意）
}
```

### 数値の意味

| フィールド | 単位 | 実用範囲 | 備考 |
|-----------|------|---------|------|
| CastTime | game tick (1/20秒) | 現状全アクティブが `40`（2秒） | 0だと即時発動 |
| CoolTime | game tick | 現状全アクティブが `20`（1秒） | 連続使用防止 |
| MPCost | MP値 | 4〜60程度 | MaxMP = Soul（最大100） |

### マスターデータ

**ファイル:** `maf-command-generator/savedata/active/`（配下の `*.json` を `JsonStore` がマージロードする）

24種のアクティブが定義されている。IDは `{名前}{番号}` 形式（例: `prominence01`, `healing01`）。

---

## 2. 生成パイプライン

```
savedata/active/
    ↓ ActiveEntity.Load()
model/active/Active[]
    ↓ export.BuildActiveArtifacts()
[]ActiveEffectFunction { ID, Body, Book, SpellBody }
    ↓ WriteActiveArtifacts()        ↓ WriteActiveDebugArtifacts()    ↓ WriteActiveSpellArtifacts()
generated/active/effect/{id}.mcfunction    generated/active/give/{id}.mcfunction    generated/active/spell/{id}.mcfunction
```

### 生成される成果物（各アクティブにつき3ファイル）

1. **effect/{id}.mcfunction** — スペル効果スクリプト。`Script[]` をそのまま結合したもの
2. **give/{id}.mcfunction** — デバッグ用。`give @p {本のNBT} 1` コマンド
3. **spell/{id}.mcfunction** — 実行時ロード用。最新の `cost/cast/cooltime/title/description` を `oh_my_dat:...maf.magic.casting` に書き込む

### NBT変換の流れ

**ファイル:** `maf-command-generator/app/domain/export/convert/active.go`

```
ActiveToBook(entry)
  → activeSpellBookModel(entry)  // spellBookModel構築
    → spellCustomData(entry)        // {maf:{active_id:"..."}}
  → .ToGiveItem()                   // minecraft:book[...] 形式に組み立て

ActiveCastingDataSNBT(entry)
  → generated/active/spell/{id}.mcfunction の casting データ
```

---

## 3. アイテムとしてのNBT構造

アクティブは `minecraft:book` アイテムとして以下のコンポーネントを持つ:

```snbt
minecraft:book[
  minecraft:item_name={text:"ヒーリング"},
  minecraft:lore=[{text:"周囲に即時回復+リジェネ"},{text:"消費MP:13 詠唱時間:40"}],
  minecraft:consumable={consume_seconds:99999,animation:"bow",has_consume_particles:false},
  minecraft:custom_data={
    maf:{
      active_id:"healing01"
    }
  }
]
```

### custom_data の設計意図

- **active_id**: アクティブの一意識別子。実行時に `generated/active/spell/{id}` を呼び出すために使う
- **cost/cast/cooltime/title/description はアイテムに埋め込まない**。再エクスポート後、既存IDの詠唱時間やMP消費を実行時ロード関数から読むため

### consumable の設計意図

- `consume_seconds:99999` — 事実上消費できない（アドバンスメントで使用を検知するためだけに consumable が必要）
- `animation:"bow"` — 弓を引くモーションで詠唱を表現
- `has_consume_particles:false` — 食べ物パーティクルを抑制

---

## 4. ランタイム実行フロー（Datapack側）

詳細は **magic-casting スキル** を参照。ここではアクティブ固有の部分のみ:

### トリガー

**ファイル:** `datapacks/.../advancement/use_active.json`

プレイヤーが本を右クリック（using_item）→ アドバンスメント付与 → `maf:magic/use_active` を呼び出し

### use_active.mcfunction の処理

1. アドバンスメントを revoke（再トリガー可能に）
2. `mafCastTime <= -1` チェック（詠唱中でないか）
3. `mafCoolTime <= 0` チェック（クールダウン中でないか）
4. `oh_my_dat` ストレージを準備し、古い casting データを削除
5. 手持ちアイテムの `maf.active_id` を `maf:magic/exec/load_active_spell` に渡す
6. `$function maf:generated/active/spell/$(id)` が `oh_my_dat:...maf.magic.casting` に最新データを書き込む
7. casting データが作成されたら `maf:magic/exec/set_magic` を呼び出し（→ magic-casting スキル参照）

### エフェクトディスパッチ

**ファイル:** `datapacks/.../magic/cast/run_active_effect.mcfunction`

```mcfunction
$function maf:generated/active/effect/$(id)
```

`cast/exec.mcfunction` が `kind:"active"` を確認した後、マクロ展開で `id` から対応する生成済み関数を呼び出す。

---

## 5. スペル効果スクリプトのパターン

`Script[]` に記述するmcfunctionコマンドの典型パターン:

```mcfunction
# ダメージ系: undead は逆効果になるため分岐
execute as @e[distance=1..8,type=#maf:undead] run effect give @s minecraft:instant_health 1 1
execute as @e[distance=1..8,type=!#maf:undead] run effect give @s minecraft:instant_damage 1 1

# ブロック操作系
fill ~-7 ~-7 ~-7 ~7 ~7 ~7 fire replace air

# 味方への回復/バフ
execute as @a[distance=..10] run effect give @s minecraft:regeneration 10 1

# 演出: サウンド + テルロー
playsound minecraft:entity.blaze.shoot master @a ~ ~ ~ 1.0 0.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は プロミネンス を唱えた！"}]
```

### 視線処理の亜種: sight 円範囲ターゲット

視線先の一点（例: 4ブロック先）を中心に円範囲で敵を拾う場合は、`maf:common/sight/eyes_circle_tagged` を使う。

```mcfunction
function maf:common/sight/eyes_circle_tagged {forward:4.0,radius:3.5,particleCount:80}
effect give @e[type=#maf:enemymob,tag=maf_sight_circle_target] minecraft:poison 10 10
```

- 引数は `forward` / `radius` / `particleCount` の3つが必須
- 対象タグは `maf_sight_circle_target`
- 共通処理側で `minecraft:witch` パーティクルを出すため、個別スペル演出は重複しないように調整する
- マクロ関数内で `$(...)` を使う行は、行頭に `$` が必須（例: `$execute ... $(radius) ...`）

### 使用するエンティティタグ

| タグ | 意味 |
|------|------|
| `#maf:undead` | アンデッド分類（ゾンビ、スケルトン等） |
| `#maf:enemymob` | 敵モブ全般 |
| `#maf:friendmob` | 味方モブ（パーティ） |

---

## 6. アイテムとのバインド

アイテムの `ItemMaf.ActiveID` にアクティブIDをセットすると、そのアイテムにアクティブがバインドされる。

**ファイル:** `maf-command-generator/app/domain/model/item/types.go`

```go
type ItemMaf struct {
    ActiveID  string `json:"activeId,omitempty"`
    PassiveID   string `json:"passiveId,omitempty"`
    PassiveSlot int    `json:"passiveSlot,omitempty"`
}
```

バインドされたアイテムは、エクスポート時に `maf.active_id` と `minecraft:consumable` が付与され、アクティブ本と同じ実行時ロード経路で使用できる。

---

## 7. ドロップシステムとの連携

アクティブは `DropRef` を通じてモンスタードロップや宝箱から入手できる。

```go
DropRef{Kind: "active", RefID: "healing01", Weight: 10, CountMin: 1, CountMax: 1}
```

- `Kind: "active"` の場合、`Slot` は設定不可（バリデーションエラー）
- `RefID` は存在するアクティブIDでなければならない（`DBMaster.HasActive()` で検証）

---

## 8. アクティブ追加手順

1. `savedata/active/ai_workspace.json` など、`savedata/active/` 配下の `*.json` にエントリ追加
2. `Script` にmcfunctionコマンドを記述
3. `make run/export` で `generated/active/effect/{id}.mcfunction`、`give/{id}.mcfunction`、`spell/{id}.mcfunction` を生成
4. ゲーム内で `/function maf:generated/active/give/{id}` でテスト
