# Repository Guidelines

## Project Goal
For the near term, this repository’s primary goal is to migrate `../tools` to Go.

## Workflow
Prefer running project commands via `make` targets when a suitable target exists.

## Lint / Staticcheck

Lint は `make lint` を使う。

`go vet` / `staticcheck` を直接叩く必要がある場合は、キャッシュ先を明示して実行する。

```bash
GOCACHE=/tmp/maf-command-editor-go-cache go vet ./...
GOCACHE=/tmp/maf-command-editor-go-cache XDG_CACHE_HOME=/tmp/maf-command-editor-cache go tool staticcheck ./...
```

Run all check commands: ```make check```

## minecraft server 

Minecraft local server is run with docker, ../compose.yml .

## mcstacker

入力画面は mcstacker を参考にする。https://mcstacker.net/?cmd=give


## old projects info

以前のプロジェクト情報が必要なら /tmp で git clone をするなどして内容を確認して。本プロジェクトの /tmp にもある。

project1: https://github.com/para7/Minecraft_Datapack
project2: https://github.com/para7/magic_and_frontier
