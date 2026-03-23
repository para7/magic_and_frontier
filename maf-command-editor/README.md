# maf-command-editor

## CLI の実行方法 (`run-cli`)

`Makefile` の `run-cli` ターゲットから CLI を実行できます。  
引数は `ARGS` で渡します。

```bash
make run-cli ARGS='validate'
make run-cli ARGS='export'
```

`run-cli` の実体は次のコマンドです。

```bash
go run ./app/cmd/tools2-cli $(ARGS)
```

利用できるコマンド:

- `validate`: savedata と export 設定のバリデーション
- `export`: バリデーション後に datapack を出力
