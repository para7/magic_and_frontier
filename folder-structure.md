# フォルダ構成と役割

## リポジトリ全体

```
minecraft-docker/
├── datapacks/          # Minecraft データパック本体
├── maf-command-generator/  # データパック生成ツール（Go製）
├── docs/               # ドキュメント
├── tmp/                # 一時ファイル・調査用クローンなど
└── compose.yml         # Docker Compose（Minecraftサーバー起動）
```

---

## datapacks/

Minecraft に読み込ませるデータパックを格納するディレクトリ。

### datapacks/magic_and_frontier/

メインのデータパック。魔法と冒険をテーマにしたゲームシステムを実装する。

```
data/maf/function/
├── load/           # ロード時の初期化処理（bootstrap等）
├── system/         # tick毎の基本処理・プレイヤーID管理・スコア管理
├── magic/          # 魔法システム（詠唱・MP管理・エフェクト発動）
│   ├── cast/       # 詠唱処理（実行・キャンセル・tick管理）
│   ├── mp/         # MP表示・管理
│   └── exec/       # 魔法効果の実行
├── skill/          # スキルシステム（魔法以外のアクティブスキル）
├── passive/        # パッシブスキル処理
├── soul/           # ソウル関連システム
├── enemy/          # 敵AI・敵スキル処理
├── generated/      # maf-command-generator が自動生成するファイルの出力先
│   ├── enemy/
│   ├── grimoire/
│   └── passive/
├── devtools/       # 開発用コマンド（reinstall、サンプル本等）
└── test/           # テスト用関数

data/maf/advancement/   # アドバンスメント定義
data/maf/loot_table/    # ルートテーブル定義
data/maf/tags/          # エンティティ・ブロック・アイテムタグ定義
data/minecraft/tags/    # Minecraftのタグ上書き（関数タグ等）
```

### datapacks/devtool/

開発専用パック。ゲームルール設定・デバッグ用リセット処理などを提供する。本番環境には含めない。

```
data/dev/function/
├── gamerule.mcfunction  # ゲームルール一括設定
├── load.mcfunction
└── reset.mcfunction     # ワールドリセット処理
```

### datapacks/sample_pack/

新しいデータパックを作るときのテンプレート。最小構成のサンプルコード。

---

## maf-command-generator/

`magic_and_frontier` データパックのコマンド生成ツール。Go で実装されており、`savedata/` に定義したゲームデータをもとに、`generated/` フォルダへ `.mcfunction` ファイルを自動出力する。

```
maf-command-generator/
├── app/
│   ├── cli/            # CLIエントリーポイント（editor / validate / export コマンド）
│   ├── domain/
│   │   ├── model/      # ドメインモデル定義
│   │   │   ├── enemy/          # 敵モデル
│   │   │   ├── enemyskill/     # 敵スキルモデル
│   │   │   ├── grimoire/       # 魔導書（魔法アイテム）モデル
│   │   │   ├── item/           # アイテムモデル
│   │   │   ├── loottable/      # ルートテーブルモデル
│   │   │   ├── passive/        # パッシブスキルモデル
│   │   │   ├── spawntable/     # スポーンテーブルモデル
│   │   │   └── treasure/       # トレジャーモデル
│   │   ├── export/     # データパック出力処理
│   │   │   └── convert/        # モデル → mcfunction 変換ロジック
│   │   ├── master/     # マスターデータ管理（ID採番など）
│   │   └── custom_validator/   # カスタムバリデーションルール
│   ├── files/          # ファイルI/O・JSON読み書きユーティリティ
│   └── minecraft/      # Minecraftルートテーブル生成ユーティリティ
├── savedata/           # ゲームデータ定義ファイル（JSON）
│   ├── enemy.json          # 敵定義
│   ├── enemy_skill.json    # 敵スキル定義
│   ├── grimoire.json       # 魔導書定義
│   ├── item.json           # アイテム定義
│   ├── passive.json        # パッシブスキル定義
│   ├── spawn_table.json    # スポーンテーブル定義
│   ├── loottables.json     # ルートテーブル定義
│   ├── treasure.json       # トレジャー定義
│   ├── id_counters.json    # IDの採番状態管理
│   └── spawn_tables/       # スポーンテーブルの詳細データ
├── config/
│   └── export_settings.json  # エクスポート先パスなどの設定
├── minecraft/          # Minecraft本体データ（バージョン別の参照情報）
│   └── 1.21.11/
└── bin/                # ビルド済みバイナリ（mcg, mce等）
```

### 主要コマンド

| コマンド | 内容 |
|---------|------|
| `mcg editor` | Web UIエディタを起動 |
| `mcg validate` | savedata のバリデーション |
| `mcg export` | バリデーション後にデータパックを出力 |

---

## tmp/

一時ファイル置き場。旧リポジトリのクローン・設計メモ・調査ログなどを格納する。バージョン管理対象外。

---

## docs/

プロジェクトのドキュメント置き場。

| ファイル | 内容 |
|---------|------|
| `folder-structure.md` | 本ドキュメント（フォルダ構成と役割） |
| `minecraft-knowledge-log.md` | Minecraftコマンド仕様の調査ログ |
| `prompt.md` | エージェント向けプロンプトメモ |
