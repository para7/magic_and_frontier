---
paths:
  - "app/domain/model/**/types.go"
---

# エンティティ型定義

`types.go` は JSON タグと validate タグ付きの構造体を定義する。

## 現在のエンティティ

| エンティティ | パッケージ | 主要フィールド |
|-------------|-----------|--------------|
| Grimoire | `model/grimoire` | ID, CastID, CastTime, MPCost, Script, Title, Description |
| Item | `model/item` | ID, ItemID, SkillID, CustomName, Lore, Enchantments, NBT 等 |
| Passive | `model/passive` | ID, Name, SkillType(sword/bow/axe), Description, Script |
| EnemySkill | `model/enemyskill` | ID, Name, Description, Script |
| Enemy | `model/enemy` | ID, MobType, Name, HP, Equipment, EnemySkillIDs, DropMode, Drops |
| SpawnTable | `model/spawntable` | ID, SourceMobType, Dimension, MinX~MaxZ, BaseMobWeight, Replacements |
| Treasure | `model/treasure` | ID, TablePath, LootPools |
| LootTable | `model/loottable` | ID, Memo, LootPools |

## 共有型（`model/types.go`）

- `DropRef`: アイテム・グリモア・バニラアイテムへの参照 + ドロップ設定（Kind, RefID, Weight, CountMin, CountMax）
- `EquipmentSlot`: エネミー装備スロット（Kind, RefID, Count, DropChance）
- `Equipment`: 6スロット（Mainhand, Offhand, Head, Chest, Legs, Feet）
- `ReplacementEntry`: スポーンテーブルのモブ差替エントリ（EnemyID, Weight）

## 新規エンティティ追加手順

1. `app/domain/model/<entity>/types.go` に型定義（json/validate タグ込み）
2. `app/domain/model/<entity>/entity.go` に `MafEntity` 実装
3. `app/domain/model/interfaces.go` の `DBMaster` に `Has*` を追加
4. `app/domain/master/master.go` にフィールド追加、`NewDBMaster` で `Load()`、`Has*` 実装
5. 必要に応じて `app/domain/export/interfaces.go` と export 実装を拡張
6. `savedata/*.json` とテストを追加し、`make check` を通す
