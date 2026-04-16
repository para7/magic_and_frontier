---
name: grimoire
description: グリモア（魔導書）システムの設計リファレンス。グリモアのデータモデル・NBT構造・生成パイプライン・ランタイム実行フローを網羅する。グリモアの新規追加、エフェクト修正、アイテムへのバインド、スペルのデバッグなど、グリモアに関わる作業全般で参照すること。grimoire, 魔導書, スペル, 魔法書, 呪文, spell, cast などのキーワードが出たら使う。
---

# グリモアシステム

グリモアはプレイヤーが使う「魔導書」アイテム。本（`minecraft:book`）として配布され、右クリックで詠唱→効果発動する。

## 関連スキル

- **magic-casting**: 詠唱パイプライン（キャスト・クールダウン・MP消費の共通基盤）。グリモアとパッシブの両方が使う
- **passive**: パッシブスキルシステム。詠唱パイプラインを共有するがトリガー方式が異なる
- **ohmydat**: プレイヤー個別ストレージ。詠唱データの一時保存先
- **maf-export**: Go ジェネレータのエクスポートパイプライン全体

---

## 1. データモデル（Generator側）

### Grimoire 構造体

**ファイル:** `maf-command-generator/app/domain/model/grimoire/types.go`

```go
type Grimoire struct {
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
| CastTime | game tick (1/20秒) | 現状全グリモアが `40`（2秒） | 0だと即時発動 |
| CoolTime | game tick | 現状全グリモアが `20`（1秒） | 連続使用防止 |
| MPCost | MP値 | 4〜60程度 | MaxMP = Soul（最大100） |

### マスターデータ

**ファイル:** `maf-command-generator/savedata/grimoire.json`

24種のグリモアが定義されている。IDは `{名前}{番号}` 形式（例: `prominence01`, `healing01`）。

---

## 2. 生成パイプライン

```
savedata/grimoire.json
    ↓ GrimoireEntity.Load()
model/grimoire/Grimoire[]
    ↓ export.BuildGrimoireArtifacts()
[]GrimoireEffectFunction { ID, Body, Book }
    ↓ WriteGrimoireArtifacts()        ↓ WriteGrimoireDebugArtifacts()
generated/grimoire/effect/{id}.mcfunction    generated/grimoire/give/{id}.mcfunction
```

### 生成される成果物（各グリモアにつき2ファイル）

1. **effect/{id}.mcfunction** — スペル効果スクリプト。`Script[]` をそのまま結合したもの
2. **give/{id}.mcfunction** — デバッグ用。`give @p {本のNBT} 1` コマンド

### NBT変換の流れ

**ファイル:** `maf-command-generator/app/domain/export/convert/grimoire.go`

```
GrimoireToBook(entry)
  → grimoireSpellBookModel(entry)  // spellBookModel構築
    → spellCustomData(entry)        // {maf:{grimoire_id:...,spell:{...}}}
      → grimoireSpellFragment(entry) // spell:{kind:"grimoire",...}
  → .ToGiveItem()                   // minecraft:book[...] 形式に組み立て
```

---

## 3. アイテムとしてのNBT構造

グリモアは `minecraft:book` アイテムとして以下のコンポーネントを持つ:

```snbt
minecraft:book[
  minecraft:item_name={text:"ヒーリング"},
  minecraft:lore=[{text:"周囲に即時回復+リジェネ"},{text:"消費MP:13 詠唱時間:40"}],
  minecraft:consumable={consume_seconds:99999,animation:"bow",has_consume_particles:false},
  minecraft:custom_data={
    maf:{
      grimoire_id:"healing01",
      spell:{
        kind:"grimoire",
        id:"healing01",
        cost:13,
        cast:40,
        cooltime:20,
        title:"ヒーリング",
        description:"周囲に即時回復+リジェネ"
      }
    }
  }
]
```

### custom_data の設計意図

- **grimoire_id**: グリモアの一意識別子。アイテム検索用
- **spell**: 詠唱パイプラインが読み取る統一フォーマット。`kind` フィールドで grimoire/passive を区別する
  - **cost/cast/cooltime**: ランタイムでスコアボードにロードされる数値
  - **title**: MPバー表示に使用
  - **description**: 現在ランタイムでは未使用（将来の詠唱UI用）

### consumable の設計意図

- `consume_seconds:99999` — 事実上消費できない（アドバンスメントで使用を検知するためだけに consumable が必要）
- `animation:"bow"` — 弓を引くモーションで詠唱を表現
- `has_consume_particles:false` — 食べ物パーティクルを抑制

---

## 4. ランタイム実行フロー（Datapack側）

詳細は **magic-casting スキル** を参照。ここではグリモア固有の部分のみ:

### トリガー

**ファイル:** `datapacks/.../advancement/use_grimoire.json`

プレイヤーが本を右クリック（using_item）→ アドバンスメント付与 → `maf:magic/use_grimoire` を呼び出し

### use_grimoire.mcfunction の処理

1. アドバンスメントを revoke（再トリガー可能に）
2. `mafCastTime <= -1` チェック（詠唱中でないか）
3. `mafCoolTime <= 0` チェック（クールダウン中でないか）
4. `oh_my_dat` ストレージを準備し、古い casting データを削除
5. 手持ちアイテムに `maf.spell` があるか確認
6. `maf.spell` を `oh_my_dat: ...maf.magic.casting` にコピー
7. `maf:magic/exec/set_magic` を呼び出し（→ magic-casting スキル参照）

### エフェクトディスパッチ

**ファイル:** `datapacks/.../magic/cast/run_grimoire_effect.mcfunction`

```mcfunction
$function maf:generated/grimoire/effect/$(id)
```

`cast/exec.mcfunction` が `kind:"grimoire"` を確認した後、マクロ展開で `id` から対応する生成済み関数を呼び出す。

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

### 使用するエンティティタグ

| タグ | 意味 |
|------|------|
| `#maf:undead` | アンデッド分類（ゾンビ、スケルトン等） |
| `#maf:enemymob` | 敵モブ全般 |
| `#maf:friendmob` | 味方モブ（パーティ） |

---

## 6. アイテムとのバインド

アイテムの `ItemMaf.GrimoireID` にグリモアIDをセットすると、そのアイテムにグリモアがバインドされる。

**ファイル:** `maf-command-generator/app/domain/model/item/types.go`

```go
type ItemMaf struct {
    GrimoireID  string `json:"grimoireId,omitempty"`
    PassiveID   string `json:"passiveId,omitempty"`
    PassiveSlot int    `json:"passiveSlot,omitempty"`
}
```

バインドされたアイテムは、エクスポート時に `spell` フラグメントが `custom_data` に埋め込まれ、グリモア本と同じように使用できる。

---

## 7. ドロップシステムとの連携

グリモアは `DropRef` を通じてモンスタードロップや宝箱から入手できる。

```go
DropRef{Kind: "grimoire", RefID: "healing01", Weight: 10, CountMin: 1, CountMax: 1}
```

- `Kind: "grimoire"` の場合、`Slot` は設定不可（バリデーションエラー）
- `RefID` は存在するグリモアIDでなければならない（`DBMaster.HasGrimoire()` で検証）

---

## 8. グリモア追加手順

1. `savedata/grimoire.json` にエントリ追加
2. `Script` にmcfunctionコマンドを記述
3. `make run/export` で `generated/grimoire/effect/{id}.mcfunction` と `give/{id}.mcfunction` を生成
4. ゲーム内で `/function maf:generated/grimoire/give/{id}` でテスト
