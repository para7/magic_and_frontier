---
paths:
  - "app/domain/master/**/*.go"
---

# master 層の規約

`master.DBMaster` は全 `MafEntity` を保持し、model 層と export 層の両方の `DBMaster` インターフェースを実装するハブ。

## 責務

- `NewDBMaster(cfg)` で全エンティティを初期化・`Load()` 済みの状態で返す
- `Has*` メソッド群: `model.DBMaster` インターフェースの実装（リレーションバリデーション用）
- `Get*ByID` / `List*` メソッド群: `export.DBMaster` インターフェースの実装 + CLI 用の追加アクセサ
- `ValidateAll()`: 全エンティティの `ValidateAll` + SpawnTable 座標範囲重複チェック
- `List*` はスライスの defensive copy を返す

## エンティティ追加時

1. フィールドに `model.MafEntity[xxx.Xxx]` を追加
2. `NewDBMaster` 内で初期化と `Load()` を追加
3. `Has*` メソッドを実装
4. 必要に応じて `Get*ByID` / `List*` を追加
5. `ValidateAll()` にエンティティを追加
