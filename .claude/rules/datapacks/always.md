---
paths: datapacks/magic_and_frontier
---

# datapacks

Minecraft データパック群。3つのパックと設計メモで構成される。

Minecraft Version 26.1

## magic_and_frontier/ — メインデータパック

名前空間 `maf`。RPG風の魔法システムを実装するデータパック。

### data/maf/function/ — 手書きファンクション

- `load.mcfunction` / `tick.mcfunction` — エントリポイント（vanilla の load/tick タグから呼ばれる）
- `entered_world.mcfunction` — ワールド参加時処理

#### magic/ — 魔法コアシステム
- `load.mcfunction` / `tick.mcfunction` — 魔法システムの初期化とティック
- `use_grimoire.mcfunction` — 魔法書使用時のメイン処理
- `player_init.mcfunction` — プレイヤー初期化
- `setdb.mcfunction` — oh_my_dat ストレージ操作
- `cast/` — 詠唱パイプライン
  - `exec.mcfunction` — 詠唱実行
  - `cancel.mcfunction` — 詠唱キャンセル
  - `tick.mcfunction` — 詠唱中ティック処理
  - `run_grimoire_effect.mcfunction` — 魔法書エフェクト実行ディスパッチ
  - `run_passive_apply.mcfunction` — パッシブ適用ディスパッチ
- `mp/` — MP管理
  - `mp_manage.mcfunction` — MP消費・回復ロジック
  - `mpbar.mcfunction` / `mpbar_init.mcfunction` — MPバー表示
  - `mpbar_per_player_dispatch.mcfunction` / `mpbar_player_macro.mcfunction` — プレイヤーごとのMPバー
- `passive/` — パッシブ効果
  - `tick.mcfunction` — パッシブティック
  - `run_effect.mcfunction` — 装備パッシブ効果実行
  - `run_mainhand_effect.mcfunction` — メインハンド装備パッシブ
  - `run_bow_effect.mcfunction` — 弓パッシブ効果実行
  - `tag_passive_arrow.mcfunction` — パッシブ矢タグ付け
  - `on_arrow_hit.mcfunction` — 矢着弾時処理
- `exec/set_magic.mcfunction` — 魔法セット

#### system/ — スコアボード・ID管理
- `score/prescore.mcfunction` / `score/afterscore.mcfunction` — スコアボードの前処理/後処理
- `set_player_id/run.mcfunction` / `set_player_id/do_not_call.mcfunction` — プレイヤーID割り当て

#### その他
- `devtools/reinstall.mcfunction` — 再インストール
- `devtools/passive_clear.mcfunction` — パッシブクリア
- `enemy/tick.mcfunction` — 敵ティック処理
- `skill/sword_slash.mcfunction` — 剣スキル
- `soul/tick.mcfunction` — ソウルシステムティック
- `test/tick.mcfunction` — テスト用ティック

### data/maf/function/generated/ — 自動生成ファンクション

maf-command-generator で生成。直接編集禁止。

- `grimoire/effect/` — 魔法書エフェクト（24種: sweep, teleport, healing, tornado, barrier, etc.）
- `grimoire/give/` — 魔法書giveコマンド（上記と対応する24種）
- `passive/effect/` — パッシブ効果（passive_1, regeneration, test_bow_passive）
- `passive/give/` — パッシブ装備giveコマンド
- `passive/apply/` — パッシブ装備適用処理
- `passive/bow/` — 弓パッシブ効果
- `enemy/spawn/` — 敵スポーン（poison_zombie, drop_test, drop_test_bow）
- `enemy/skill/` — 敵スキル（main, near_poison）

### data/maf/advancement/ — 進捗（イベントトリガー）
- `use_grimoire.json` — 魔法書使用検知（consumable + using_item）
- `entered_world.json` — ワールド参加検知
- `arrow_hit.json` — 矢着弾検知

### data/maf/loot_table/generated/ — 自動生成ルートテーブル
- `enemy/loot/` — 敵ドロップテーブル
- `item/items_1.json` — アイテムプール

## p7BaseSystem/ — 前提パック

名前空間 `p7b`。他のデータパックが依存する共通基盤。

### data/p7b/function/ — ユーティリティ関数
- `load.mcfunction` / `tick.mcfunction` / `init.mcfunction` — エントリポイント
- `sword.mcfunction` — 剣関連処理
- `warp.mcfunction` — ワープ
- `generate_rand.mcfunction` — 乱数生成
- `killme.mcfunction` — 自殺コマンド

### data/p7b/tags/ — 共有タグ定義
- `entity_type/` — mob分類（mobs, undead, zombies, skeletons, spiders, water_enemy, enemymob, enemymob_notboss, friendmob）
- `block/` — ブロック分類（air, water）
- `item/` — アイテム分類（swords）

## sample_pack/ — テスト用パック

Docker マウント検証用の最小パック。`sample:ping` を load 時に実行するだけ。
