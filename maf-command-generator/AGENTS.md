# maf-command-generator

Minecraft データパック向けゲームコンテンツ（魔導書・アイテム・エネミー・スキル等）をJSONで管理し、バリデーション後にデータパック（`.mcfunction` / loot table JSON）としてエクスポートするGoアプリケーション。

開発初期段階のため、過去互換より設計改善を優先する。

## スタック

- Go 1.25.0（go.mod）
- `github.com/go-playground/validator/v10`（バリデーション）
- `github.com/a-h/templ`（将来のUI用、現在未使用）

## コマンド

- `make check` — フル検証（generate + tidy + format + lint + build + test）
- `make run/validate` — バリデーション実行
- `make run/export` — バリデーション + エクスポート実行

## ディレクトリ概要

- `app/cli/` — editor, export, validate のエントリポイント
- `app/domain/model/` — エンティティ定義・CRUD・バリデーション（`MafEntity[T]` パターン）
- `app/domain/export/` — データパック出力（model に直接依存しない）
- `app/domain/master/` — 全 model を束ねるハブ（`DBMaster`）
- `app/files/` — JSON I/O ユーティリティ
- `app/minecraft/` — バニラ loot table 読み込み
- `savedata/` — エンティティ JSON データ
- `config/` — エクスポート設定

## 注意事項

- アーキテクチャ詳細は `.claude/rules/` にパスベースのルールとして配置済み
