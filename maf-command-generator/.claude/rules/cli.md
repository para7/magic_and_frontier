---
paths:
  - "app/cli/**/*.go"
  - "app/main.go"
---

# CLI 層の規約

## main.go

- サブコマンド分岐のみ（`editor` / `validate` / `export`）
- 引数なしはエラー終了（usage 表示）

## cli パッケージ

- `Editor` / `Validate` / `Export` は `int` を返し、`main` が `os.Exit` する
- `Export` は必ず `Validate` を先に実行し、失敗時は出力処理を中断する
- 表示フォーマット調整は `cli` 側で行い、ドメイン層へ持ち込まない
- `editor` は現在スタブ（将来 templ ベースの Web UI を予定）
