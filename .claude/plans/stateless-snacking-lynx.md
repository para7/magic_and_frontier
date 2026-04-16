# JsonStore ディレクトリ分割対応

## Context

item, enemy 等のエンティティを大量追加すると単一JSONファイルが肥大化する。`JsonStore[T]` をディレクトリ対応に拡張し、手動でファイル分割できるようにする。

## 方針

- `JsonStore[T].Path` がファイルなら従来動作、ディレクトリなら `*.json` を全読み込み・マージ
- Load時に各エントリの元ファイル名を `sourceMap` (map[string]string: ID→basename) で保持
- Save時は `sourceMap` に基づき元ファイルへ書き戻し。新規エントリは `entity.json` へ
- 変更は `files` 層のみ。model/export 層は無変更

## 変更ファイル

### 1. `app/files/json.go` — メイン実装

**フィールド追加:**
```go
type JsonStore[T any] struct {
    Path      string
    isDir     bool
    sourceMap map[string]string // id -> basename
}
```

**ID抽出ヘルパー:**
```go
type idHolder struct { ID string `json:"id"` }
type rawEntriesFile struct { Entries []json.RawMessage `json:"entries"` }
func extractID[T any](entry T) (string, error) { ... }
```

**Load:**
- `os.Stat` でファイル/ディレクトリ判定
- ファイル → 既存 `loadFile` ヘルパー
- ディレクトリ → `loadDir`: glob `*.json`, 各ファイル読み込み, rawEntriesFile でID抽出, sourceMap構築, マージ
- 重複ID検出 → エラー返却

**Save:**
- `isDir == false` → 既存 `saveFile` ヘルパー
- `isDir == true` → `saveDir`: extractIDでグループ化, sourceMap参照で振り分け, 新規は `entity.json`, 空になったファイルは `{"entries":[]}` で残す

### 2. `app/files/json_test.go` — テスト追加

- ファイルモード: 既存動作維持
- ディレクトリモード: 複数ファイルマージ、sourceMap正確性
- Save後の振り分け: 元ファイルへの書き戻し、新規→entity.json
- 重複ID検出
- 空ディレクトリ

### 3. `app/files/config.go` — 今回は変更なし

パス変更は手動。ユーザーが分割したいエンティティのみ:
```
"savedata/item.json" → "savedata/item"  (ディレクトリ)
```

## マイグレーション手順

1. `mkdir savedata/item`
2. `mv savedata/item.json savedata/item/entity.json`
3. config.go のパスを `"savedata/item"` に変更
4. ユーザーが手動でファイル分割

## 実装順序

1. 既存 Load/Save → `loadFile`/`saveFile` ヘルパーに切り出し
2. `idHolder`, `rawEntriesFile`, `extractID` 追加
3. `isDir`, `sourceMap` フィールド追加
4. `loadDir` 実装 + Load にディスパッチ追加
5. `saveDir` 実装 + Save にディスパッチ追加
6. テスト追加
7. `make check`

## 検証

- `make check` (generate + tidy + format + lint + build + test)
- 既存の単一ファイルモードで `make run/validate` が通ること
