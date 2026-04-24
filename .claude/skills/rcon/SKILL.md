---
name: rcon
description: How to run commands in minecraft server.
---

## マイクラサーバーへのコマンド発行方法

作業ディレクトリ: `/home/para7/workspaces/minecraft-docker`

### 基本コマンド

```bash
make mc-cmd CMD='<minecraftコマンド>'
```

例:
```bash
make mc-cmd CMD='list'
make mc-cmd CMD='reload'
make mc-cmd CMD='function maf:debug/test'
```

### 仕組み

- `minecraft` コンテナと `backup` コンテナの両方が起動している必要がある
- `backup` コンテナ内の `rconclt` ツールを使って RCON 接続
- RCON パスワードとポートは `minecraft` コンテナの `/mc/data/server.properties` から自動取得

### トラブルシューティング

コンテナ内シェルで直接確認したい場合:
```bash
make mc-shell
```
