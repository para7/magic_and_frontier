# maf_command_editor

## コマンド

`mcg` は次の 3 コマンドを持ちます。

- `editor`: Web サーバー起動
- `validate`: savedata と export 設定のバリデーション
- `export`: バリデーション後に datapack を出力

引数なしの `mcg` は `mcg editor` と同じです。

## `make` での実行

```bash
make run                       # editor 起動
make run-cmd ARGS='validate'   # validate 実行
make run-cmd ARGS='export'     # export 実行
```

実体コマンド:

```bash
go run ./app editor
go run ./app validate
go run ./app export
```
