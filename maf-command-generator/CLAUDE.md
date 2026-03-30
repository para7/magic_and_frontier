# maf-command-generator

Minecraft データパック向けゲームコンテンツ（魔導書・アイテム・エネミー・スキル等）をJSONで管理し、バリデーション後にデータパック（`.mcfunction` / loot table JSON）としてエクスポートするGoアプリケーション。

本システムは内部で使うPC専用ツールです。スマホ向けのcssは不要。

## プロジェクト設計概要

- app/cli/* : エントリポイント
- app/domain/model/* : savedata のJSONと紐づくエンティティ、MafEntity の実装。
- app/domain/export/* : エクスポート用実装。

## コマンド

- `make check` — フル検証（generate + tidy + format + lint + build + test）
- `make run/validate` — バリデーション実行
- `make run/export` — バリデーション + エクスポート実行

## 旧プロジェクト

移行元の参照が必要な場合:

- project1: https://github.com/para7/Minecraft_Datapack
- project2: https://github.com/para7/magic_and_frontier
- 移行前のプロジェクト: "../maf-command-editor" 
