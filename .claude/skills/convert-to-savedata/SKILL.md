---
name: convert-to-savedata
description: datapacks/magic_and_frontier/data/maf/function/devtools/tmp*.mcfunction の内容を maf-command-generator/savedata/{item|enemy}/convert.json に変換して取り込むスキル。tmp の give/summon コマンドを手動でエディタに写したくない、アイテム/エネミーの savedata に反映したい、convert.json へ追記したい、といった要求で使う。convert, savedata 化, tmp 取り込み, give→item, summon→enemy などのキーワードで呼ぶ。
disable-model-invocation: true
---

# tmp.mcfunction → savedata 変換

`datapacks/.../devtools/tmpN.mcfunction` の `give`/`summon` コマンドを読み取り、`maf-command-generator/savedata/{item|enemy}/convert.json` に新規エントリとして書き込む。エクスポートして自動生成コマンドが tmp と一致するかを比較スクリプトで確認する。

## 引数

`/convert-to-savedata <tmpName> <type>`

- `tmpName` — `tmp` または `tmp2` など、`datapacks/magic_and_frontier/data/maf/function/devtools/{tmpName}.mcfunction` の basename。
- `type` — `item` または `enemy` のみ。それ以外はエラー終了。

## 入出力パス

| type | 入力 tmp | 出力 savedata | 検証用自動生成先 |
|------|---------|---------------|----------------|
| item | `datapacks/magic_and_frontier/data/maf/function/devtools/{tmpName}.mcfunction` | `maf-command-generator/savedata/item/convert.json` | `datapacks/magic_and_frontier/data/maf/function/generated/item/give/{id}.mcfunction` |
| enemy | 同上 | `maf-command-generator/savedata/enemy/convert.json` | `datapacks/magic_and_frontier/data/maf/function/generated/enemy/spawn/{id}.mcfunction` |

## 実行フロー

1. **tmp 読込**: 指定 tmp を読む。先頭コメント/空行はスキップし、最初の非空行を対象コマンドとして扱う。
2. **コマンド分解**: `give @p <itemId>[components] <count>` または `/summon <mobType> ~ ~ ~ {NBT}` を分解。詳細は下記「変換ルール」。
3. **id 決定**: `YYYYMMDDHHMM` 形式のタイムスタンプ (例 `202604241530`) を id とする。重複したらサフィックス `_1`, `_2` を足す。ユーザーに後で修正させる前提で、マイクラで有効な文字 (小文字英数 + `_`) のみを使う。
4. **convert.json 読込/新規作成**: 対象パスに `convert.json` が無ければ `{"entries": []}` で初期化。既存なら読込。
5. **エントリ追記**: `entries` の末尾に新エントリを追加してファイル保存。
6. **export 実行**: プロジェクトルートで `make run/export` を実行。RCON の reload ステップが失敗しても export 自体が成功していれば続行。バリデーションエラーの扱いは下記「バリデーション失敗時の自動修正」を参照。
7. **比較**: `scripts/compare_tmp_vs_generated.py` で tmp と自動生成ファイルを比較。不足があれば修正→再 export。
8. **完了報告**: ユーザーへ以下を必ず報告:
   - 追加した id と savedata のパス
   - 自動修正した項目があれば「何を・どう推測して・どの値にしたか」
   - 比較スクリプトの結果 (OK or 残存差分)

## 変換ルール (item)

`give @p <itemId>[<components>] <count>` → `{id, itemId, maf:{}, minecraft:{components:{...}}}`

- `itemId` に `minecraft:` が付いていなければ付与して savedata の root `itemId` に入れる。
- components 内の各キーに `minecraft:` プレフィックスが無ければ付与し、`minecraft.components` 配下に置く。
- 文字列値の component は JSON テキストコンポーネントに包む:
  - `custom_name="foo"` → `{"text":"foo"}`
  - `lore=["a","b"]` → `[{"text":"a"},{"text":"b"}]`
- `enchantments={"sharpness":10}` は中のキーにも `minecraft:` を付ける: `{"minecraft:sharpness":10}`
- `attribute_modifiers=[{...}]` は配列をそのまま。中の `amount/operation/slot/type/id` はそのまま保持。
- 型サフィックス `1b/1f/1.0d` などは数値として解釈し、JSON では素の数値で書く。

### 例

tmp 入力:
```
give @p bundle[custom_name="ビジネスバッグ",lore=["お世話になっております。"],enchantments={"sharpness":10}] 1
```

savedata 出力:
```json
{
  "id": "202604241530",
  "itemId": "minecraft:bundle",
  "maf": {},
  "minecraft": {
    "components": {
      "minecraft:custom_name": {"text": "ビジネスバッグ"},
      "minecraft:lore": [{"text": "お世話になっております。"}],
      "minecraft:enchantments": {"minecraft:sharpness": 10}
    }
  }
}
```

## 変換ルール (enemy)

`/summon <mobType> ~ ~ ~ {<NBT>}` → `{id, mobType, memo:"", minecraft:{...(ほぼ全NBT)}, maf:{...デフォルト}}`

- `mobType` に `minecraft:` が無ければ付与。
- **原則: NBT のトップレベルフィールドは原則そのまま `minecraft.*` に積む。** export 側 (`maf-command-generator/app/domain/export/convert/enemy.go`) は `entry.Minecraft` を deepCopy してから maf 側の Tags / equipment / Passengers をマージするので、`minecraft` 側に置いたフィールドは保存される (上書きされるのは後述の数件のみ)。
- 型サフィックスは除去し素の数値/真偽値に変換する (`1b` → `1`, `0b` → `0`, `1f` → `1.0`)。
- `Tags` はメインエネミーでは `minecraft.Tags` に残す。export が `maf_enemy*` / `EnemySkill` 系を追加マージする (既存 Tags は保持される)。
- `Passengers` は savedata トップレベル `passengers:[...]` に移す (各 Passenger も同じルールで再帰)。Passenger だけ `Tags` は `maf.tags` 側に入れる規約。
- `equipment:{<slot>:{id,count}}` / `drop_chances` はそのまま `minecraft.equipment` / `minecraft.drop_chances` に置く。`maf.equipment` は空 `{}` にする (maf 側は savedata のアイテム参照 (`kind: item`) を使いたい時用で、NBT から機械的に判定できないため)。
- `data:{...}` (custom_data 相当) も `minecraft.data` に残す。

### `maf` のデフォルト値

NBT から導けないフィールドは以下で埋め、修正レポートに「デフォルトを投入した」と明記する:

- `maf.equipment = {}`
- `maf.enemySkillIds = []` (既存 Tags から逆引きはしない)
- `maf.dropMode = "replace"`
- `maf.drops = []`
- `memo = ""`

### 例

tmp2 入力 (抜粋):
```
/summon creeper ~ ~ ~ {CustomName:"name",Health:1f,equipment:{feet:{id:"minecraft:stone",count:1}},drop_chances:{feet:0.0},attributes:[{id:"minecraft:max_health",base:1}],Tags:["cutstomTag"]}
```

savedata 出力:
```json
{
  "id": "202604241530",
  "mobType": "minecraft:creeper",
  "memo": "",
  "minecraft": {
    "CustomName": "name",
    "Health": 1.0,
    "equipment": {"feet": {"id": "minecraft:stone", "count": 1}},
    "drop_chances": {"feet": 0.0},
    "attributes": [{"id": "minecraft:max_health", "base": 1}],
    "Tags": ["cutstomTag"]
  },
  "maf": {
    "equipment": {},
    "enemySkillIds": [],
    "dropMode": "replace",
    "drops": []
  }
}
```

### エクスポートが必ず上書きするフィールド (比較スクリプトが差分として拾うが無視してよい)

- `DeathLootTable` → 常に `maf:generated/enemy/loot/<id>` に書き換え。tmp に `"minecraft:empty"` 等が書かれていても生成物では別値になる。
- `Tags` に `maf_enemy` / `maf_enemy_<id>` / `maf_vh_checked` / (maf.enemySkillIds がある場合) `EnemySkill` / `<skillId>` / `maf_enemy_skill_<skillId>` が必ず追加される。
- 既存 Tags は保持される (重複しない場合)。

## バリデーション失敗時の自動修正

`make run/export` が validate で落ちた場合は、エラー出力を読んで推測で修正してからリトライする。代表的な失敗と対処:

| エラー | 推測した原因 | 自動修正 |
|--------|-------------|---------|
| `id must match ...` | id に不正文字 | タイムスタンプを再生成 or 英数+アンダースコアに正規化 |
| `maf_slug_id` | enemy id の形式違反 | 小文字化・記号除去 |
| `trimmed_required: itemId` | itemId 未設定 or 空 | tmp から再抽出して埋める |
| `dropMode oneof=append replace` | 値が別物 | `replace` に固定 |
| 参照先エンティティが存在しない (grimoireId/passiveId/bowId 等) | tmp に `data:{maf:{...}}` が入っていた | 当該 maf フィールドを空に戻す |

修正するたびに「何を」「どのキーを見てそう推測したか」「どの値にしたか」を記録し、最終報告で必ず提示する。推測ができない/自信が無い修正は自動で行わずユーザーに判断を仰ぐ。

## 比較スクリプト

`scripts/compare_tmp_vs_generated.py` は tmp の `give`/`summon` と自動生成の同種コマンドを SNBT パースして構造比較する。「tmp にある key/value が generated にも存在するか」だけをチェック (generated 側に追加で `maf_enemy_*` Tags や `DeathLootTable` などが増えるのは許容)。

```bash
python3 .claude/skills/convert-to-savedata/scripts/compare_tmp_vs_generated.py \
  datapacks/magic_and_frontier/data/maf/function/devtools/{tmpName}.mcfunction \
  datapacks/magic_and_frontier/data/maf/function/generated/{item/give|enemy/spawn}/{id}.mcfunction
```

- 終了コード `0`: 不足なし (完了条件)。
- 終了コード `1`: 不足あり。stdout に `path: 内容` 形式で列挙される。
- 終了コード `2`: パース失敗。

### 差分の読み方

- `key: missing in generated` → savedata の書き漏れ。tmp からコピーし直す。
- `[i]: no matching element in generated` → 配列要素が対応する項目を持たない。attributes や attribute_modifiers の取りこぼしが多い。
- `value mismatch` → 値が食い違っている。型変換ミスか文字列→テキストコンポーネント化漏れの可能性。

## 完了条件

- 比較スクリプトが終了コード 0 を返す。
- 返さない場合は savedata を修正→ `make run/export` →再比較、をユーザーに状況を報告した上で数回ループ。差分が埋まらなければユーザーへ判断を仰ぐ。
