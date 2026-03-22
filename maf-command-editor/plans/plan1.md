 # Dataflow Gap Closure: Full Schema + Export Overhaul

  ## Summary

  現状の CRUD ベース実装を、dataflow-user.md の設計に合わせて一括で作り直す。対象は保存スキーマ、ID 採番、Web/API 入力、export、検証、テスト一式。

  実装は 3 系統で並行化する。

  1. 保存スキーマ・ID 採番・アプリ検証
  2. Web/UI・HTTP API・フォーム変換
  3. export・生成物テスト

  互換変換は入れない。旧 savedata/*.json はそのまま読めても新スキーマ検証で止め、必要なら手修正する。

  ## Key Changes

  ### ID / 保存

  - 新規カウンタ state を savedata/id-counters.json に追加し、items, grimoire, skill, enemyskill, enemy, treasure, castid を一元管理する。
  - 新規 ID は items_<n>, grimoire_<n>, skill_<n>, enemyskill_<n>, enemy_<n>, treasure_<n> で採番する。
  - castid は整数の連番で採番し、同じカウンタ state で管理する。
  - 既存ドメインの uuid_any 前提をやめ、保存 ID は「非空文字列 + 各 prefix 形式」に切り替える。

  ### ドメイン / API / UI

  - GrimoireEntry を id, castid, castTime, mpCost, title, description, script, updatedAt に置換する。variants は削除する。
  - SkillEntry を id, name?, description?, script, updatedAt に置換する。itemId は削除する。
  - ItemEntry に optional skillId を追加する。1 item につき skill は 0 または 1 件。
  - EnemySkillEntry を id, name?, description?, script, updatedAt に置換する。cooldown, trigger は削除する。
  - TreasureEntry を id, mode, tablePath, lootPools, updatedAt に置換する。
      - mode は custom | override
      - tablePath は export 先の loot table path
      - lootPools.kind は minecraft_item | item | grimoire
      - custom の tablePath は state 内で重複禁止
  - EnemyEntry を id, mobType, name, hp, attack, defense, moveSpeed, equipment, enemySkillIds, dropMode, drops, updatedAt に置換する。
      - dropMode は append | replace
      - drops は enemy 自身が直接保持する
      - dropTableId と treasure 参照は削除する
      - equipment は少なくとも mainhand, offhand, head, chest, legs, feet を構造化して持つ
  - Web/UI は新規作成時にサーバ側で ID / castid を払い出し、編集画面では read-only 表示にする。
  - Skill 画面から item 選択を外し、Item 画面で optional skill を選ぶ形に反転する。
  - Grimoire 画面から variants 入力を外し、cast time と MP cost の 2 項目にする。
  - Treasure 画面に mode ラジオと tablePath 入力を追加する。
  - Enemy 画面に mobType, equipment, dropMode, drops を追加し、dropTableId を除去する。

  ### Export

  - Grimoire は 1 エントリにつき以下を生成する。
      - 実行本体 function: spellFunctionDir/<grimoire_id>.mcfunction
      - 本アイテム loot table: spellLootDir/<grimoire_id>.json
      - 中央 dispatcher: spellFunctionDir/selectexec.mcfunction
  - Grimoire の本体 function には script をそのまま出力する。book 側 custom data には grimoire_id, castid, cast, cost を入れ、raw script は埋め込まない。
  - Dispatcher は castid ごとに execute if score ... run function <namespace>:.../<grimoire_id> を自動生成する。
  - Item export は skillId がある場合だけ maf_skill:1b と maf_skill_id:"skill_<n>" を custom data に付与する。
  - Skill / EnemySkill export は ID そのままのファイル名で mcfunction を生成する。
  - Treasure export は mode/tablePath に従って保存先を決める。
      - custom: 指定 path に新規 table を出力
      - override: 指定 path をそのまま上書き
  - Enemy export は enemy_id ベースの固定 loot table 名を使う。
      - loot table: enemyLootDir/<enemy_id>.json
      - summon function: enemyFunctionDir/<enemy_id>.mcfunction
      - summon NBT では DeathLootTable に上記 table を設定する
      - mobType, 名前, 属性, 装備, enemy skill 用タグ/識別子を summon 側へ反映する

  ## Validation / Error Policy

  - 互換変換は行わない。旧 schema の variants, skill.itemId, enemy.dropTableId, treasure UUID-only model などは保存・export 前 validation で明示エラーにする。
  - castid は一意必須。
  - item.skillId は存在する skill のみ許可。
  - treasure.tablePath は custom 内で重複禁止、空文字禁止。
  - lootPools / enemy drops は kind ごとに参照先を検証する。
  - enemySkillIds は存在確認し、重複は正規化で除去する。

  ## Test Plan

  - ドメインテスト
      - ID / castid 採番
      - grimoire 新 schema 検証
      - item-skill 参照
      - treasure mode + tablePath 検証
      - enemy direct drop / equipment / dropMode 検証
      - 旧 schema 入力が validation error になること
  - HTTP / Web テスト
      - 新フォーム保存
      - skill 選択が item 側に移ること
      - grimoire variants 廃止
      - treasure custom/override path 重複エラー
      - enemy dropTableId 廃止後の保存
  - Export テスト
      - grimoire script function がそのまま出ること
      - selectexec.mcfunction が castid 分岐を持つこと
      - item に maf_skill / maf_skill_id が付くこと
      - treasure が mode/tablePath に従って正しい場所へ出ること
      - enemy summon が mobType と DeathLootTable を反映すること
      - skill / enemy skill が ID 名ファイルで出ること
  - 統合テスト
      - 全 state を新 schema で保存し go test ./... が通ること
      - /save と POST /api/save で新生成物一式が揃うこと

  ## Assumptions

  - 互換性は不要。旧データは自動移行しない。
  - item の mcstacker 完全再実装は今回やらず、現行ハイブリッド入力を拡張する。
  - grimoire は内部では独立エンティティのまま持ち、loot/input では専用参照を許可する。
  - enemy は treasure を参照せず、自身で drops を持つ。
  - enemy の一般ドロップは inline ではなく DeathLootTable 参照で出す。