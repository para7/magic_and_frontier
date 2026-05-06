---
name: buff-system
description: maf datapack の独自バフシステム設計リファレンス。buff, バフ, custom buff, 状態効果, tick_function, destructor_function, parry, パリィ, resistance など、汎用バフ処理やパリィの挙動確認・修正・追加で参照すること。
---

# 独自バフシステム

`maf` 名前空間のプレイヤー別一時効果ランタイム。各プレイヤーの oh_my_dat private storage にバフ配列を保持し、毎 tick tick 関数を実行しながら残り tick を減算する。期限切れ時には destructor 関数を呼ぶ。

## 関連ファイル

- `datapacks/magic_and_frontier/data/maf/function/tick.mcfunction`: `execute as @a at @s run function maf:buff/tick`
- `datapacks/magic_and_frontier/data/maf/function/common/buff/set.mcfunction`: バフ追加・同一 `buff_id` refresh の共通入口
- `datapacks/magic_and_frontier/data/maf/function/buff/add.mcfunction`: `maf.buff_entry` を `maf.buff` 配列へ append
- `datapacks/magic_and_frontier/data/maf/function/buff/tick.mcfunction`: バフ配列を queue 化して処理開始
- `datapacks/magic_and_frontier/data/maf/function/buff/process_queue.mcfunction`: 1件ずつ tick 実行、残り tick 更新、期限切れ判定
- `datapacks/magic_and_frontier/data/maf/function/buff/run_tick.mcfunction`: `$function maf:buff/tick/$(tick_function)`
- `datapacks/magic_and_frontier/data/maf/function/buff/run_destructor.mcfunction`: `$function maf:buff/destructor/$(destructor_function)`
- `datapacks/magic_and_frontier/data/maf/function/buff/tick/*.mcfunction`: バフ種別ごとの毎 tick 処理
- `datapacks/magic_and_frontier/data/maf/function/buff/destructor/*.mcfunction`: バフ種別ごとの終了処理

## Storage 形

対象プレイヤーで `function #oh_my_dat:please` 実行後、以下を使う。

```mcfunction
oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff
```

`maf.buff` は配列。各 entry は現状この形。

```snbt
{
  buff_id:"parry",
  tick:10,
  tick_function:"parry_tick",
  destructor_function:"parry_destructor"
}
```

処理中の一時領域:

- `maf.buff_entry`: `common/buff/set` が作る追加予定エントリ
- `maf.buff_queue`: `buff/tick` が当該 tick の処理対象としてコピーするキュー
- `maf.buff_current`: `process_queue` が処理中の1件

## 追加・更新フロー

基本入口:

```mcfunction
function maf:common/buff/set {buff_id:"parry",tick:10,tick_function:"parry_tick",destructor_function:"parry_destructor"}
```

`common/buff/set` の挙動:

1. 対象プレイヤーの oh_my_dat storage を開く。
2. `maf.buff` がなければ空配列で初期化する。
3. 同一 `buff_id` の既存エントリがあれば `maf.buff_current` にコピーし、既存エントリの `destructor_function` を実行する。
4. 同一 `buff_id` の既存エントリを `maf.buff` から削除する。
5. 新しい entry を `maf.buff_entry` に作成し、`buff/add` で `maf.buff` に append する。

この仕様では同一 `buff_id` は refresh 扱いで、同時多重スタックしない。

## Tick フロー

`tick.mcfunction` から全プレイヤーに対して `buff/tick` を呼ぶ。

`buff/tick`:

1. `maf.buff[0]` がなければ終了。
2. `maf.buff` を `maf.buff_queue` にコピーする。
3. `maf.buff` を空配列にする。
4. `buff/process_queue` を再帰的に呼び、queue を1件ずつ処理する。
5. 最後に `maf.buff_queue` と `maf.buff_current` を削除する。

`buff/process_queue`:

1. queue 先頭を `maf.buff_current` に移し、queue から削除する。
2. `tick_function` があれば `maf:buff/tick/{tick_function}` を macro function で呼ぶ。
3. `maf.buff_current.tick` を scoreboard `tmp` に読み、1減算して storage に戻す。
4. `tick <= 0` なら `destructor_function` を呼ぶ。
5. `tick >= 1` なら更新後の `maf.buff_current` を `maf.buff` に append する。
6. 次の queue 要素へ進む。

## パリィ実装

入力データ:

- `maf-command-generator/savedata/grimoire/ai_workspace.json` の `parry01`
- 生成物: `datapacks/magic_and_frontier/data/maf/function/generated/grimoire/effect/parry01.mcfunction`

発動時:

```mcfunction
function maf:common/buff/set {buff_id:"parry",tick:10,tick_function:"parry_tick",destructor_function:"parry_destructor"}
effect give @s minecraft:resistance 1 9 true
```

毎 tick:

- `buff/tick/parry_tick.mcfunction`
- cloud particle を出す。
- `HurtTime:9s` を検出したらパリィ音を鳴らす。

終了時:

- `buff/destructor/parry_destructor.mcfunction`
- `effect clear @s minecraft:resistance`

## 注意点

- `destructor_function` は終了時だけでなく、同一 `buff_id` の refresh 時にも呼ばれる。
- destructor は副作用を最小化すること。現状の `parry_destructor` は `minecraft:resistance` を全消去するため、他システム由来の resistance と競合しうる。
- `tick_function` / `destructor_function` は macro 経由で関数名に展開される。外部入力化する場合は許可済み関数名だけを渡すこと。
- `process_queue` は `tmp` scoreboard を使う。tick 関数内で `tmp` を使う場合、残り tick 減算前後の値を壊さないようにする。
- `buff/tick` 実行中に追加されたバフは空にした `maf.buff` 側へ append されるため、次 tick から処理される。
