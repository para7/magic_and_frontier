---
paths: ["datapacks/magic_and_frontier"]
---

# datapacks

Minecraft データパック群。Minecraft Version 26.1

## magic_and_frontier/ — メインデータパック

名前空間 `maf`。RPG 風の魔法システム（アクティブ・パッシブ・弓パッシブ・エネミー・ソウル）を実装。

### ディレクトリ構造

```
data/maf/
├── function/
│   ├── load.mcfunction / tick.mcfunction   ← エントリポイント
│   ├── entered_world.mcfunction
│   ├── magic/                  ← 魔法コアシステム
│   │   ├── cast/               ← 詠唱パイプライン（exec/cancel/tick/run_*）
│   │   ├── mp/                 ← MP 管理・MP バー表示
│   │   ├── use_active.mcfunction
│   │   ├── player_init.mcfunction
│   │   └── setdb.mcfunction    ← oh_my_dat ストレージ操作
│   ├── passive/                ← パッシブランタイム
│   │   ├── tick.mcfunction     ← スロット1〜3 + メインハンド + 弓着弾
│   │   ├── run_effect.mcfunction / run_bow_effect.mcfunction
│   │   ├── on_arrow_hit / on_bow_hit / on_melee_hit  ← advancement コールバック
│   │   └── tag_passive_arrow.mcfunction
│   ├── bow/                    ← 弓パッシブ共有ランタイム
│   │   ├── tag_bow_arrow / prepare_hit_arrow / on_bow_hit
│   │   ├── process_hit_arrows / resolve_hit_arrow / run_bow_effect
│   │   └── tick_flying / run_flying / tick_ground / run_ground
│   ├── enemy/                  ← エネミーティック・スポーン置換エントリ
│   │   ├── tick.mcfunction     ← EnemySkill 駆動 + maf_vh_checked 未付与を replace へ
│   │   └── replace.mcfunction  ← → generated/enemy/replace/main へディスパッチ
│   ├── system/score/           ← スコアボード前後処理・プレイヤーID 割当
│   ├── soul/tick.mcfunction    ← ソウルシステム
│   ├── skill/sword_slash.mcfunction
│   ├── devtools/               ← reinstall, passive_clear
│   ├── warp / generate_rand / sword / killme  ← ユーティリティ
│   └── generated/              ← maf-command-generator 生成（直接編集禁止）
│       ├── active/{effect,give}/
│       ├── item/give/
│       ├── passive/{effect,give,apply,bow}/
│       ├── bow/{flying,ground}/
│       ├── enemy/{spawn,skill,replace}/
│       └── （内容は export_settings.json で定義）
├── advancement/                ← イベントトリガー
│   ├── use_active / entered_world
│   ├── arrow_hit → passive/on_arrow_hit
│   └── melee_hit → passive/on_melee_hit
├── loot_table/generated/       ← 自動生成ルートテーブル
│   ├── enemy/loot/
│   └── item/
└── tags/                       ← entity_type / block / item 分類
```

### mcfunction スクリプト規約

- 一時計算には `tmp` / `tmp2` スコアボードを使用（`load.mcfunction` で定義済み）
- `mafMP` 等の用途固定スコアボードを一時計算に流用しない

## sample_pack/

Docker マウント検証用の最小パック。`sample:ping` を load 時に実行するだけ。

## 変更後チェック（必須）

- `datapacks/` 配下のファイルを編集した場合は、RCON 経由で必ず `reload` を実行して再読込結果を確認する
- 実行コマンド: `make mc-cmd CMD='reload'`
- `reload` 実行結果にエラー（構文エラー、関数ロード失敗など）が出ないことを確認する
