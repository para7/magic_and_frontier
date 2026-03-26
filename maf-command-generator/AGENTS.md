# project info

移行前のプロジェクト "../maf-command-editor" を最小構成で書き直し中。

# Structure

- app/cli/ -- editor, export, validate の3機能のエントリポイント。
- app/domain/model/ -- 各 JSON に対応した、統一規格でのモデリングと抽象化、バリデーションなど。
- app/domain/export/ -- データパックの出力処理。
- app/domain/master/ -- 全 model を束ねて model 間連携を扱う。
- app/files/ -- ファイル操作系ユーティリティ。

# 重要クラス

## MafEntity

app/domain/model で実装している、各 JSON 操作やバリデーションの実装。

## DBMaster

あらゆる機能へアクセスできるハブとなる存在。データ間連携を取り持つ。
cli, editor などは DBMaster を薄くラッピングする形で機能を提供する。

# Repository Guidelines

- build check command: ```make check```

