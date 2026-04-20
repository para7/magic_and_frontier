---
name: enemy-replace
description: バニラモブをカスタムエネミーに重み付き確率で置換するスポーンテーブルシステム（enemy-replace）の設計リファレンス。SpawnTable データモデル、座標・距離・ディメンション・mobType による条件分岐、baseMob による残存モブの属性上書き、`maf_vh_checked` タグと `maf_vh_rand` スコアを用いたランタイムフローを網羅する。スポーンテーブルの新規追加、置換確率の調整、バニラモブ属性上書き、座標範囲の変更、`enemy/replace` / `enemy/tick` の挙動調査などで参照すること。spawn_table, spawntable, enemy replace, エネミー置換, モブ置換, baseMob, ベースモブ, replacement, maf_vh_rand, maf_vh_checked, sourceMobType, minDistance, spawnTable エクスポート などのキーワードで使う。
---

# enemy-replace / SpawnTable システム

バニラの自然スポーンモブ（zombie 等）を、特定のディメンション・距離・座標範囲内で一定確率にカスタムエネミーへ置き換える仕組み。マスターデータ（SpawnTable）から mcfunction ディスパッチャ群を生成し、各エンティティが 1 度だけチェックされて置換 or 残存（必要なら属性上書き）を決める。

## 関連スキル

- **maf-export**: 全エクスポートパイプラインの位置付け。`BuildSpawnTableArtifacts` はそのひとつ
- **passive** / **bow-passive** / **grimoire**: 他のエクスポート対象。共通のビルド/ライト分離パターン

---

## 1. データモデル（SpawnTable）

**ファイル:** `maf-command-generator/app/domain/model/spawntable/types.go`

```go
type SpawnTable struct {
    ID            string
    SourceMobType string                   // 例: "minecraft:zombie"
    Dimension     string                   // overworld / the_nether / the_end
    MinDistance   int                      // ワールド原点(0,100,0)からの距離 [0, 99999999]
    MaxDistance   int
    MinX, MaxX    int                      // [-99999999, 99999999]
    MinY, MaxY    int
    MinZ, MaxZ    int
    BaseMob       *BaseMob                 // nil 可。nil の場合バニラ残存枠なし
    Replacements  []model.ReplacementEntry // min=1、EnemyID + Weight
}

type BaseMob struct {
    Weight     int                // バニラモブが残る抽選枠（0 なら 100% 置換）
    Attributes *BaseMobAttributes // nil 可
}

type BaseMobAttributes struct {
    HP, Attack, Defense, MoveSpeed *float64 // 全て optional（omitempty）
}
```

### アクセサと後方互換

- `SpawnTable.GetBaseMobWeight()` / `GetBaseMobAttributes()` — BaseMob が nil でも安全に取り出す
- 旧 savedata の `"baseMobWeight": N`（フラットスカラー）は `UnmarshalJSON` で `BaseMob{Weight: N}` に昇格される。既存マスターデータを壊さないためのブリッジなので、新規データは常に `baseMob: { weight, attributes }` ネスト形式で書く

### Replacement エントリ

```go
type ReplacementEntry struct {
    EnemyID string // maf:generated/enemy/spawn/<EnemyID> を呼ぶ
    Weight  int    // 抽選重み
}
```

- `EnemyID` は Enemy エンティティの ID と一致。`master.HasEnemy` で存在確認される
- 合計重み（BaseMob.Weight + Σ Replacement.Weight）は必ず > 0

---

## 2. savedata フォーマット

`savedata/spawn_table/*.json` は他エンティティと違い、ファイル単位の `coordinates` を共有する独自形式:

```json
{
  "coordinates": {
    "dimension": "minecraft:overworld",
    "minDistance": 0,
    "maxDistance": 1000,
    "minX": -99999999, "maxX": 99999999,
    "minY": -64,       "maxY": 320,
    "minZ": -99999999, "maxZ": 99999999
  },
  "entries": [
    {
      "sourceMobType": "minecraft:zombie",
      "baseMob": { "weight": 30, "attributes": { "hp": 100 } },
      "replacements": [
        { "enemyId": "drop_test",     "weight": 30 },
        { "enemyId": "poison_zombie", "weight": 40 }
      ]
    }
  ]
}
```

### ロード規則（`entity.go` の `Load()`）

- `savedata/spawn_table/` 配下の `*.json` を `filepath.Glob` で走査
- 各ファイルの `coordinates` を `entries[i]` に複製してから `SpawnTable` に展開（エントリ単位で上書きしたい場合はエントリ側に同名フィールドを置けば優先される）
- `id` が空のエントリは `buildSpawnTableID(fileBase, index, sourceMobType)` で自動付与
  - 小文字化、namespace 除去（`minecraft:zombie` → `zombie`）、非英数は `_` にまとめる、末尾の `_` 除去
  - 例: `over1.json` + `minecraft:zombie` + index=0 → `over1_zombie_1`
- `JsonStore[T]` を使わないので、新規エンティティでこの形式を真似る場合は独自パーサを書く必要がある

---

## 3. バリデーション

`entity.go` の `ValidateRelation` で以下をチェック:

| 項目 | 内容 |
|------|------|
| `minDistance <= maxDistance` | 距離区間の整合 |
| `minX/Y/Z <= maxX/Y/Z` | 座標区間の整合 |
| `Replacements[i].EnemyID` 存在 | `mas.HasEnemy` |
| `Replacements` 内の EnemyID 重複禁止 | `seen` マップ |
| 合計重み > 0 | `BaseMob.Weight + Σ Replacement.Weight` |
| 属性範囲 | `validateBaseMobAttributes`: HP 1..100000, Attack/Defense/MoveSpeed 0..100000 |

`master.ValidateAll` 側では更に `RangesOverlap`（下記）で、同一 `SourceMobType` × `Dimension` の SpawnTable 同士の座標＋距離範囲重複を検出する。

### RangesOverlap（`overlap.go`）

距離 (`MinDistance..MaxDistance`) と X/Y/Z の 4 区間すべてが重なったときに「重複」と判定する AND 条件。ひとつでも分離していれば OK。新規テーブルの座標を考える際は、既存テーブルのどれか一軸でずらせば衝突しない。

---

## 4. エクスポートパイプライン

**ファイル:** `maf-command-generator/app/domain/export/spawntable.go`

```go
BuildSpawnTableArtifacts(dmas, replaceLogicalDir, enemyLogicalDir) []SpawnTableArtifact
WriteSpawnTableArtifacts(dir, artifacts) error

type SpawnTableArtifact struct {
    ID        string // テーブル ID
    Table     string // {id}.mcfunction の本文
    MainEntry string // main.mcfunction の 1 行（ディスパッチャ条件）
}
```

### 出力先

- 論理ディレクトリ: `export_settings.json` の `spawnTable` キー（デフォルト `generated/enemy/replace`）
- 物理パス: `{outputRoot}/data/maf/function/{spawnTable}/`
- ファイル: `main.mcfunction` + テーブル数ぶんの `{id}.mcfunction`

### main.mcfunction（ディスパッチャ）

各テーブルにつき以下の 1 行を連結:

```
execute if dimension <dim> \
        if entity @s[x=0,y=100,z=0,distance=<minD>..<maxD>] \
        if entity @s[x=<minX>,y=<minY>,z=<minZ>,dx=<dx>,dy=<dy>,dz=<dz>,type=<srcMob>] \
        run function maf:<spawnTable>/<id>
```

- `distance` はワールド原点 `(0,100,0)` 基準。Y=100 固定は高さを距離計算から事実上無視するための便法
- `dx/dy/dz` は `max - min`（min 側を起点とする AABB）

### {id}.mcfunction（個別置換）

```
execute store result score @s maf_vh_rand run random value 1..<totalWeight>
execute if score @s maf_vh_rand matches <start>..<end> run function maf:<enemy>/<enemyId>   # 各 Replacement
...
execute if score @s maf_vh_rand matches 1..<replaceWeight>              run function maf:killme
execute if score @s maf_vh_rand matches <replaceStart>..<totalWeight>   run data merge entity @s <baseMobNBT>   # 属性ありのとき
```

- 置換レンジ `1..replaceWeight` に当たった場合のみ `killme` でバニラモブ個体を消す。カスタムエネミーは別個体として summon される設計
- `replaceWeight + 1..totalWeight` はバニラモブ残存枠。`BaseMob.Attributes` が指定されている場合のみ `data merge entity` で属性を書き込む（属性が全部 nil ならそのまま放置）

### baseMobMergeNBT

`BaseMobAttributes` を Minecraft の属性 NBT に変換する:

| フィールド | 出力 |
|-----------|------|
| `HP` | `Health:<v>f` と `Attributes[{Name:generic.max_health, Base:<v>}]` の両方を書く |
| `Attack` | `Attributes[{Name:generic.attack_damage, Base:<v>}]` |
| `Defense` | `Attributes[{Name:generic.armor, Base:<v>}]` |
| `MoveSpeed` | `Attributes[{Name:generic.movement_speed, Base:<v>}]` |

- 全フィールド nil のときは `ok=false` を返し、`data merge` 行そのものを出力しない
- HP のみ `Health` と `max_health` の二重書きなのは、`max_health` を上げただけでは現在 HP は初期値のままになるため

---

## 5. ランタイムフロー（datapack 側）

### エントリポイント

- `data/maf/function/enemy/tick.mcfunction`
  ```
  execute as @e[tag=EnemySkill] at @s run function maf:generated/enemy/skill/main
  execute as @e[tag=!maf_vh_checked] at @s run function maf:enemy/replace
  ```
  `maf_vh_checked` が未付与のエンティティに対して 1 度だけ `enemy/replace` を呼ぶ

- `data/maf/function/enemy/replace.mcfunction`
  ```
  tag @s add maf_vh_checked
  function maf:generated/enemy/replace/main
  ```
  最初にタグを付けるので、置換判定で `killme` された個体も、生き残ったバニラも、summon された別エンティティも二度とここを通らない

### スコアボード

- `maf_vh_rand` (`dummy`): `load.mcfunction` で登録。`{id}.mcfunction` 内で `random value` の抽選結果を一時保持するためだけの共有レジスタ

### タグ

- `maf_vh_checked`: SpawnTable チェック済みマーカー。付いていると `tick` 側で再チェックされない

### 置換先エネミーの呼び出し

`function maf:generated/enemy/spawn/<enemyId>` は enemy エクスポート側の成果物。そちらは `@s at @s` の位置で summon するので、バニラ個体の座標に重ねて新規エンティティが湧く構造になる。その直後に同テーブルの `killme` 行でバニラ個体を void 送りにする。

---

## 6. 新しいスポーンテーブルを追加する

1. `savedata/spawn_table/<area>.json` を作成（または既存ファイルに entry を追加）
   - `coordinates` はファイル単位。同じ範囲のテーブルを複数並べるなら 1 ファイルにまとめると可読性が上がる
   - 個別 entry で座標を上書きしたい場合のみ、entry 側に `dimension` / `minX` 等を書く
2. `replacements[].enemyId` に対応する Enemy が `savedata/enemy/` に存在することを確認
3. 重複チェック: 同じ `sourceMobType` × `dimension` で、距離・座標が既存テーブルと全軸重なっていないこと（`RangesOverlap`）
4. `make run/validate` で構造 + リレーション検証
5. `make run/export` でデータパック生成
6. ゲーム内で該当ディメンション・距離・座標範囲に対象モブがスポーンした時に、意図した確率で置換/属性上書きが起きるか確認

### よくある調整

- **置換率を変えたい**: `baseMob.weight` と各 `replacements[].weight` の比を変える。合計が 100 である必要はないが、100 にしておくと比率が読みやすい
- **100% 置換にしたい**: `baseMob` を省略（nil）。それだけで `baseMobMergeNBT` も走らなくなる
- **残存バニラを強化したい**: `baseMob.attributes` に HP/Attack 等を書く。HP のみ変えたいなら `hp` だけ、それ以外は attributes 丸ごと省略して OK（`*float64` なので nil と 0 は区別される）
- **距離で難易度ゾーニング**: 同一 `sourceMobType` × `dimension` で `minDistance/maxDistance` を分けて複数テーブルを並べる。座標 AABB は全域にしておけば良い

---

## 7. 設定

**ファイル:** `maf-command-generator/config/export_settings.json`

```json
"spawnTable": "generated/enemy/replace"
```

- キー: `ExportPaths.SpawnTable`（`app/files/config.go`）
- デフォルト適用は `normalizePathOrDefault`（空欄時 `generated/enemy/replace`）
- `passive/bow` のようにハードコードされていないので、別パスに逃がしたい場合はここを書き換える

---

## 8. テスト

- `app/domain/export/enemy_export_test.go` にフィクスチャベースのテストがあり、`testdata/spawn_table/basic/` と `testdata/spawn_table/basemob_attributes/` の 2 シナリオを検証
- `basic/input/spawn_tables.json` は旧 `baseMobWeight` 形式（後方互換テスト）、`basemob_attributes/input/spawn_tables.json` は新 `baseMob.attributes` 形式
- `app/domain/model/spawntable/entity_test.go` に構造 + リレーションバリデーションの単体テスト

---

## 9. 関連ファイル一覧

### Generator（Go）

| ファイル | 役割 |
|---------|------|
| `maf-command-generator/app/domain/model/spawntable/types.go` | SpawnTable / BaseMob / BaseMobAttributes 構造体 + UnmarshalJSON |
| `maf-command-generator/app/domain/model/spawntable/entity.go` | 独自 `Load()`、`ValidateRelation`、`validateBaseMobAttributes`、`buildSpawnTableID` |
| `maf-command-generator/app/domain/model/spawntable/overlap.go` | `RangesOverlap`（距離 + X/Y/Z 4 軸 AND） |
| `maf-command-generator/app/domain/export/spawntable.go` | `BuildSpawnTableArtifacts` / `WriteSpawnTableArtifacts` / `baseMobMergeNBT` |
| `maf-command-generator/app/domain/export/interfaces.go` | `export.DBMaster.ListSpawnTables()` |
| `maf-command-generator/app/files/config.go` | `ExportPaths.SpawnTable` |
| `maf-command-generator/config/export_settings.json` | `spawnTable` キー |

### Datapack（手書き）

| ファイル | 役割 |
|---------|------|
| `data/maf/function/enemy/tick.mcfunction` | 未チェックエンティティに `enemy/replace` をディスパッチ |
| `data/maf/function/enemy/replace.mcfunction` | `maf_vh_checked` タグ付与 + 自動生成 `main` 呼び出し |
| `data/maf/function/load.mcfunction` | `maf_vh_rand` スコアボード登録 |
| `data/maf/function/killme.mcfunction` | void 送りでエンティティ削除（置換レンジで呼ばれる） |

### Datapack（自動生成）

| ファイル | 役割 |
|---------|------|
| `generated/enemy/replace/main.mcfunction` | ディスパッチャ（dimension + distance + AABB + type で分岐） |
| `generated/enemy/replace/{id}.mcfunction` | 個別テーブル。`maf_vh_rand` 抽選 → 置換エネミー呼出 / killme / 残存モブ属性上書き |

### savedata

| ファイル | 役割 |
|---------|------|
| `savedata/spawn_table/*.json` | `{coordinates, entries}` 形式のマスターデータ |
