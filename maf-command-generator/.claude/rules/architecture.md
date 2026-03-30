---
paths:
  - "app/**/*.go"
---

# レイヤーアーキテクチャ

依存方向は固定で、逆流させない。

```
main → cli → master → model / export → files, minecraft
```

- `main`: サブコマンド分岐のみ（`editor` / `validate` / `export`）
- `cli`: 入出力と終了コード管理。ドメインロジックを持たない
- `master`: モデル間連携のハブ。`model.DBMaster` と `export.DBMaster` 両方を実装
- `model`: エンティティ単位の CRUD・検証・永続化
- `export`: 変換と出力。`export.DBMaster` 経由で読むだけ
- `files`: JSON ストアと設定ロードの I/O ユーティリティ
- `minecraft`: バニラ loot table 読み込み・存在確認

# インターフェース

- `model.DBMaster`（`app/domain/model/interfaces.go`）: リレーションバリデーション用の `Has*` メソッド群
- `model.MafEntity[T]`（同上）: エンティティ共通の CRUD・検証・永続化インターフェース
- `export.DBMaster`（`app/domain/export/interfaces.go`）: エクスポート専用の読取インターフェース（`Get*`/`List*`）
- `master.DBMaster`（`app/domain/master/master.go`）: 上記2つを実装する具象型

# ビルド

- `make check` でフル検証（generate + tidy + format + lint + build + test）
