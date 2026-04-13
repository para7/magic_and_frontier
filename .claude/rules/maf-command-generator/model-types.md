---
paths:
  - "maf-command-generator/app/domain/model/**/types.go"
---

# エンティティ型定義

`types.go` は JSON タグと validate タグ付きの構造体を定義する。

## 現在のエンティティ

| エンティティ | パッケージ | 主要フィールド |
|-------------|-----------|--------------|
| Grimoire | `model/grimoire` | ID, CastTime, CoolTime, MPCost, Script, Title, Description |
| Item | `model/item` | ID, Maf(GrimoireID/PassiveID/PassiveSlot/BowID), Minecraft(ItemID/Components) |
| Passive | `model/passive` | ID, Name, Role, Condition(always/on_sword_hit), Slots, Description, Script |
| BowPassive | `model/bow` | ID, Name, Role, Slots, LifeSub, ScriptHit/Fired/Flying/Ground |
| EnemySkill | `model/enemyskill` | ID, Name, Description, Script |
| Enemy | `model/enemy` | ID, MobType, Name, HP, Equipment, EnemySkillIDs, DropMode, Drops |
| SpawnTable | `model/spawntable` | ID, SourceMobType, Dimension, MinX~MaxZ, BaseMobWeight, Replacements |
| Treasure | `model/treasure` | ID, TablePath, LootPools |
| LootTable | `model/loottable` | ID, Memo, LootPools |

## 共有型（`model/types.go`）

- `DropRef`: アイテム・グリモア・パッシブ・バニラアイテムへの参照 + ドロップ設定（Kind(minecraft_item/item/grimoire/passive), RefID, Slot, Weight, CountMin, CountMax）
- `EquipmentSlot`: エネミー装備スロット（Kind(minecraft_item/item), RefID, Count, DropChance）
- `Equipment`: 6スロット（Mainhand, Offhand, Head, Chest, Legs, Feet）
- `ReplacementEntry`: スポーンテーブルのモブ差替エントリ（EnemyID, Weight）

## 新規エンティティ追加手順

1. `maf-command-generator/app/domain/model/<entity>/types.go` に型定義（json/validate タグ込み）
2. `maf-command-generator/app/domain/model/<entity>/entity.go` に `MafEntity` 実装
3. `maf-command-generator/app/domain/model/interfaces.go` の `DBMaster` に `Has*` を追加
4. `maf-command-generator/app/domain/master/master.go` にフィールド追加、`NewDBMaster` で `Load()`、`Has*` 実装
5. 必要に応じて `maf-command-generator/app/domain/export/interfaces.go` と export 実装を拡張
6. `savedata/*.json` とテストを追加し、`make check` を通す
