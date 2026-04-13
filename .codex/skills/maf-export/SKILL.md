---
name: maf-export
description: maf-command-generator のエクスポートパイプライン設計リファレンス。Go コードがどのようにマスターデータ（JSON）からデータパック成果物（.mcfunction, loot table JSON）を生成するかを解説する。エクスポートの仕組み変更、新エンティティのエクスポート追加、生成物の構造理解、convert 層の修正などで参照すること。export, エクスポート, 生成, generate, build artifacts, convert, loot table, mcfunction生成 などのキーワードで使う。
---

# エクスポートパイプライン

Go ジェネレータ（maf-command-generator）が JSON マスターデータからデータパック成果物を生成する仕組み。

## 関連スキル

- **grimoire**: グリモア固有の変換ロジック
- **passive**: パッシブ固有の変換ロジック
- **ohmydat**: 生成されるコマンドが使う oh_my_dat ストレージ

---

## 1. アーキテクチャ全体像

```
savedata/*.json
    ↓ model/*.Entity.Load()
master.DBMaster (全エンティティを保持)
    ↓ export.ExportDatapack(dmas, config)
    ↓ (export.DBMaster インターフェース経由で読取)
Build*Artifacts() → メモリ上に成果物構築
    ↓
Write*Artifacts() → ファイル書き出し
    ↓
datapacks/magic_and_frontier/data/maf/
    ├── function/generated/  (.mcfunction)
    └── loot_table/generated/ (loot table JSON)
```

### レイヤー構造

```
main → cli → master → model / export → files, minecraft
```

- **model 層**: エンティティ単位の CRUD・バリデーション・永続化。`model.DBMaster` で Has* メソッド提供
- **export 層**: 変換と出力。`export.DBMaster` インターフェース経由で読取のみ
- **convert サブパッケージ**: 純粋変換関数群。export を逆インポートしない

---

## 2. ExportDatapack オーケストレーション

**ファイル:** `maf-command-generator/app/domain/export/export.go`

```go
func ExportDatapack(dmas DBMaster, mafconfig config.MafConfig) error
```

処理順序:

1. `config/export_settings.json` から出力パスを読み込み
2. **グリモ���**: `BuildGrimoireArtifacts` → `WriteGrimoireArtifacts` + `WriteGrimoireDebugArtifacts`
3. レガシーファイル���除（`selectexec.mcfunction`, `setup_effect_ref_map.mcfunction`）
4. **アイテム**: `BuildItemArtifacts` → `WriteItemArtifacts`
5. **パッシブ**: `BuildPassiveArtifacts` → `WritePassiveArtifacts`
6. **弓パッシブ**: `BuildBowArtifacts` → `WriteBowArtifacts`
7. **エネミ��スキル**: `BuildEnemySkillArtifacts` → `WriteEnemySkillArtifacts`
8. **エネミー**: `BuildEnemyArtifacts` → `WriteEnemyArtifacts`

### Build と Write の分離

- `Build*Artifacts()`: 副作用なし。メモリ上で成果物を構築して返す
- `Write*Artifacts()`: ファイル書き出し。Build の結果を受け取る

この分離によりテストが容易になっている。

---

## 3. 出力ディレクトリ構成

**設定:** `config/export_settings.json`

| 設定キー | デフォルト値 | 出力先 |
|---------|------------|-------|
| grimoireEffect | `generated/grimoire/effect` | `function/generated/grimoire/effect/` |
| grimoireDebug | `generated/grimoire/give` | `function/generated/grimoire/give/` |
| passiveEffect | `generated/passive/effect` | `function/generated/passive/effect/` |
| passiveBow | `generated/passive/bow` | `function/generated/passive/bow/` |
| passiveGive | `generated/passive/give` | `function/generated/passive/give/` |
| passiveApply | `generated/passive/apply` | `function/generated/passive/apply/` |
| bowFlying | `generated/bow/flying` | `function/generated/bow/flying/` |
| bowGround | `generated/bow/ground` | `function/generated/bow/ground/` |
| enemySkill | (設定必須) | `function/generated/enemy/skill/` |
| enemy | (設定必須) | `function/generated/enemy/spawn/` |
| enemyLoot | (設定必須) | `loot_table/generated/enemy/loot/` |

パスは `{outputRoot}/data/maf/{function|loot_table}/{logicalDir}/` に展開される。

---

## 4. convert サブパッケージ

**パッケージ:** `maf_command_editor/app/domain/export/convert`（import alias: `ec`）

純粋変換関数のみを持つ。`export` パッケージへの逆依存は禁止。

### ファイル一覧

| ファイル | 主要関数 | 用途 |
|---------|---------|------|
| `json.go` | `JsonString`, `ToCountValue`, `ToWeight` | SNBT文字列化、lootテーブル数値変換 |
| `grimoire.go` | `GrimoireToBook` | グリモア → 本アイテム SNBT |
| `passive.go` | `PassiveToBook` | パッシブ → 設定書アイテム SNBT |
| `book.go` | `spellBookModel.ToGiveItem` | spell本の共通組み立て |
| `item.go` | (内部) | アイテム custom_data・コンポーネント・エンチャント |
| `loottable.go` | `BuildDropLootPool`, `MergeLootTablePools` | loot pool 構築・バニラマージ |
| `enemy.go` | `ToEnemyFunctionLines` | エネミー summon コマンド生成 |

### spell フラグメント生成

`grimoire.go` と `passive.go` は共通の `spellFragment()` を使う:

```go
func spellFragment(kind, id string, slot *int, mpCost, castTime, coolTime int, title, description string) string
```

- グリモア: `slot = nil` → slot フィールドなし
- パッシブ: `slot = &N` → `slot:N` フィールドあり

### spellBookModel（共通本モデル）

```go
type spellBookModel struct {
    itemName   string   // アイテム表示名
    lore       []string // 説明文行
    customData string   // {maf:{...}} SNBT
}
```

全 spell 本に共通:
- `consumable={consume_seconds:99999,animation:"bow",has_consume_particles:false}`
- `item_name`, `lore`, `custom_data` はエンティティごとに異なる

---

## 5. エンティティ別の成果物

### グリモア

```
GrimoireEffectFunction { ID, Body, Book }
```
- **effect/{id}.mcfunction**: `Body`（Script[] の結合）
- **give/{id}.mcfunction**: `give @p {Book} 1`

### パッシブ

```
PassiveEffectFunction { ID, Body }
PassiveBowFunction    { ID, Body }
PassiveGrimoireFunction { PassiveID, Slot, FunctionID, GiveBody, ApplyBody, Book }
```
- **effect/{id}.mcfunction**: 効果本体（bow の場合は自動生成の弓検知コード）
- **bow/{id}.mcfunction**: bow パッシブの矢着弾時スクリプト
- **give/{id}_slot{N}.mcfunction**: 設定書 give コマンド
- **apply/{id}_slot{N}.mcfunction**: スロット書き込み処理

### 弓パッシブ

```
BowEffectFunction   { ID, Body }
BowHitFunction      { ID, Body }
BowFlyingFunction   { ID, Body }
BowGroundFunction   { ID, Body }
```
- **effect/bow_{id}.mcfunction**: 弓検知 + 矢タグ付け + ScriptFired（passiveEffectDir に出力）
- **bow/{id}.mcfunction**: 着弾時スクリプト（ScriptHit、passiveBowDir に出力）
- **flying/{id}_flying.mcfunction**: 飛翔中スクリプト（ScriptFlying）
- **ground/{id}_ground.mcfunction**: 着地スクリプト（ScriptGround）

### エネミースキル

```
EnemySkillFunction { ID, Body, LogicalDir }
```
- **skill/{id}.mcfunction**: スキルスクリプト

### エネミー

```
EnemyArtifact { ID, SpawnBody, LootTable }
```
- **spawn/{id}.mcfunction**: summon コマンド + 装備 + スキル設定
- **loot/{id}.json**: loot table JSON（replace/append モード対応）

---

## 6. バリデーション

エクスポート前に `make run/validate` で全エンティティのバリデーションが走る。

### バリデーション種別

1. **構造バリデーション**: validate タグによる型・範囲チェック
2. **リレーションバリデーション**: `model.DBMaster.Has*()` で参照先の存在確認
   - 例: `DropRef.Kind="grimoire"` → `HasGrimoire(RefID)` で存在チェック
3. **一括バリデーション**: 全レコードの重複IDチェック

### 設計ガードレール（再発防止）

- `ValidateDropRefs` など共通ヘルパーには汎用ルールのみ置く（存在確認・slot・範囲）
- エンティティ固有の業務制約（例: passive の `generate_grimoire` 条件）は各 `ValidateRelation` 側で判定する
- 必要な参照が増えたら `model.DBMaster` に正式メソッドを追加し、`entity.go` 内で場当たりのローカル interface を作らない
- `export/convert` で validate 相当の業務エラーを返さない。失敗は validate フェーズで早期に返す
- 生成フラグ変更で成果物が減るケースでは、`Write*Artifacts` で stale ファイル削除を実装する

### カスタムバリデータ

- `maf_slug_id`: 小文字・ハイフン・アンダースコアのみ許可
- `trimmed_required`: トリム後に空でないことを検証
- `trimmed_max=N`: トリム後の最大長チェック
- `trimmed_oneof=a b c`: トリム後の値が指定リストに含まれることを検証

---

## 7. 実行方法

```bash
make run/export    # バリデーション + エクスポート
make run/validate  # バリデーションのみ
make check         # フル検証（generate + tidy + format + lint + build + test）
```
