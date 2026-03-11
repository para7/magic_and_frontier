# Data Flow

## 1. 論理テーブル

保存元 JSON:

- `savedata/form-state.json` -> `items`
- `savedata/grimoire-state.json` -> `grimoire`
- `savedata/skill-state.json` -> `skills`
- `savedata/enemy-skill-state.json` -> `enemy_skills`
- `savedata/treasure-state.json` -> `treasures`
- `savedata/enemy-state.json` -> `enemies`

### items

```sql
CREATE TABLE items (
  id UUID PRIMARY KEY,
  minecraft_item_id TEXT,           -- Minecraft item id
);

CREATE TABLE grimoire (
  id UUID PRIMARY KEY,
  cast_id INTEGER,        -- export 用の固有値
  title TEXT,
);
CREATE TABLE skills (
  id UUID PRIMARY KEY,
  item_id UUID NOT NULL,  -- FK -> items.id
  name TEXT,,
  FOREIGN KEY (item_id) REFERENCES items(id)
);
CREATE TABLE enemy_skills (
  id UUID PRIMARY KEY,
  name TEXT,
);
-- items の配列。
CREATE TABLE treasures (
  loot_table TEXT, -- minecraft 公式に用意されている、ルートテーブル名。
  itemid UUID PRIMARY KEY, -- item_id
);
CREATE TABLE enemies (
  id UUID PRIMARY KEY,
  drop_table_id UUID NOT NULL, -- FK -> treasures.id
  name TEXT,,
  FOREIGN KEY (drop_table_id) REFERENCES treasures(id)
);
```

### treasure_loot_pools

`treasures[].lootPools[]` を展開した中間表。

```sql
CREATE TABLE treasure_loot_pools (
  treasure_id UUID NOT NULL,   -- FK -> treasures.id
  kind TEXT NOT NULL,          -- 'item' | 'grimoire'
  ref_id UUID NOT NULL,        -- kind に応じて items.id or grimoire.id
  weight NUMERIC NOT NULL,
  count_min NUMERIC NULL,
  count_max NUMERIC NULL
);
```

### enemies

```sql
```

### enemy_skill_refs

`enemies[].enemySkillIds[]` を展開した中間表。

```sql
CREATE TABLE enemy_skill_refs (
  enemy_id UUID NOT NULL,       -- FK -> enemies.id
  enemy_skill_id UUID NOT NULL, -- FK -> enemy_skills.id
  FOREIGN KEY (enemy_id) REFERENCES enemies(id),
  FOREIGN KEY (enemy_skill_id) REFERENCES enemy_skills(id)
);
```

## 2. JSON 由来の参照関係

```sql
-- skill -> item
SELECT s.id AS skill_id, i.id AS item_id
FROM skills s
JOIN items i ON i.id = s.item_id;

-- treasure loot(item)
SELECT t.id AS treasure_id, i.id AS item_id
FROM treasures t
JOIN treasure_loot_pools p ON p.treasure_id = t.id
JOIN items i ON i.id = p.ref_id
WHERE p.kind = 'item';

-- treasure loot(grimoire)
SELECT t.id AS treasure_id, g.id AS grimoire_id
FROM treasures t
JOIN treasure_loot_pools p ON p.treasure_id = t.id
JOIN grimoire g ON g.id = p.ref_id
WHERE p.kind = 'grimoire';

-- enemy -> treasure(drop table)
SELECT e.id AS enemy_id, t.id AS treasure_id
FROM enemies e
JOIN treasures t ON t.id = e.drop_table_id;

-- enemy -> enemy skill
SELECT e.id AS enemy_id, es.id AS enemy_skill_id
FROM enemies e
JOIN enemy_skill_refs r ON r.enemy_id = e.id
JOIN enemy_skills es ON es.id = r.enemy_skill_id;
```

## 3. 画面ごとのデータフロー

### `/items`

一覧表示:

```sql
SELECT id, item_id, updated_at
FROM items
ORDER BY updated_at DESC, id;
```

保存:

```sql
MERGE INTO items AS dst
USING (SELECT :id, :item_id, :updated_at) AS src(id, item_id, updated_at)
ON dst.id = src.id
WHEN MATCHED THEN UPDATE SET item_id = src.item_id, updated_at = src.updated_at
WHEN NOT MATCHED THEN INSERT (id, item_id, updated_at) VALUES (src.id, src.item_id, src.updated_at);
```

削除:

```sql
DELETE FROM items
WHERE id = :id;
```

### `/grimoire`

一覧表示:

```sql
SELECT id, cast_id, title, updated_at
FROM grimoire
ORDER BY updated_at DESC, id;
```

保存:

```sql
MERGE INTO grimoire AS dst
USING (SELECT :id, :cast_id, :title, :updated_at) AS src(id, cast_id, title, updated_at)
ON dst.id = src.id
WHEN MATCHED THEN UPDATE SET cast_id = src.cast_id, title = src.title, updated_at = src.updated_at
WHEN NOT MATCHED THEN INSERT (id, cast_id, title, updated_at) VALUES (src.id, src.cast_id, src.title, src.updated_at);
```

### `/skills`

画面表示時、`items` を候補として参照:

```sql
SELECT s.id, s.name, s.item_id, i.item_id AS item_label, s.updated_at
FROM skills s
LEFT JOIN items i ON i.id = s.item_id
ORDER BY s.updated_at DESC, s.id;

SELECT id, item_id
FROM items
ORDER BY item_id, id;
```

保存前検証:

```sql
SELECT s.id, i.id AS resolved_item_id
FROM (SELECT :id AS id, :item_id AS item_id) s
LEFT JOIN items i ON i.id = s.item_id;

-- resolved_item_id IS NULL なら validation error
```

保存:

```sql
MERGE INTO skills AS dst
USING (SELECT :id, :name, :item_id, :updated_at) AS src(id, name, item_id, updated_at)
ON dst.id = src.id
WHEN MATCHED THEN UPDATE SET name = src.name, item_id = src.item_id, updated_at = src.updated_at
WHEN NOT MATCHED THEN INSERT (id, name, item_id, updated_at) VALUES (src.id, src.name, src.item_id, src.updated_at);
```

### `/enemy-skills`

一覧表示:

```sql
SELECT id, name, updated_at
FROM enemy_skills
ORDER BY updated_at DESC, id;
```

削除前参照チェック:

```sql
SELECT e.id
FROM enemy_skill_refs r
JOIN enemies e ON e.id = r.enemy_id
WHERE r.enemy_skill_id = :id;

-- 1件でもあれば削除不可
```

### `/treasures`

画面表示時、`items` と `grimoire` を候補として参照:

```sql
SELECT t.id, t.name, t.updated_at
FROM treasures t
ORDER BY t.updated_at DESC, t.id;

SELECT id, item_id AS label
FROM items
ORDER BY item_id, id;

SELECT id, title AS label
FROM grimoire
ORDER BY title, id;
```

一覧上で loot の中身を引く:

```sql
SELECT
  t.id AS treasure_id,
  p.kind,
  p.ref_id,
  COALESCE(i.item_id, g.title) AS ref_label
FROM treasures t
JOIN treasure_loot_pools p ON p.treasure_id = t.id
LEFT JOIN items i ON p.kind = 'item' AND i.id = p.ref_id
LEFT JOIN grimoire g ON p.kind = 'grimoire' AND g.id = p.ref_id
ORDER BY t.id;
```

保存前検証:

```sql
SELECT p.*
FROM treasure_loot_pools p
LEFT JOIN items i ON p.kind = 'item' AND i.id = p.ref_id
LEFT JOIN grimoire g ON p.kind = 'grimoire' AND g.id = p.ref_id
WHERE (p.kind = 'item' AND i.id IS NULL)
   OR (p.kind = 'grimoire' AND g.id IS NULL);

-- 1件でもあれば validation error
```

保存:

```sql
MERGE INTO treasures AS dst
USING (SELECT :id, :name, :updated_at) AS src(id, name, updated_at)
ON dst.id = src.id
WHEN MATCHED THEN UPDATE SET name = src.name, updated_at = src.updated_at
WHEN NOT MATCHED THEN INSERT (id, name, updated_at) VALUES (src.id, src.name, src.updated_at);

DELETE FROM treasure_loot_pools
WHERE treasure_id = :id;

INSERT INTO treasure_loot_pools (treasure_id, kind, ref_id, weight, count_min, count_max)
VALUES (:id, :kind, :ref_id, :weight, :count_min, :count_max);
```

### `/enemies`

画面表示時、`enemy_skills` を候補として参照:

```sql
SELECT e.id, e.name, e.drop_table_id, e.updated_at
FROM enemies e
ORDER BY e.updated_at DESC, e.id;

SELECT id, name AS label
FROM enemy_skills
ORDER BY name, id;
```

敵一覧に関連情報を付ける:

```sql
SELECT
  e.id AS enemy_id,
  e.name,
  t.name AS treasure_name,
  es.id AS enemy_skill_id,
  es.name AS enemy_skill_name
FROM enemies e
LEFT JOIN treasures t ON t.id = e.drop_table_id
LEFT JOIN enemy_skill_refs r ON r.enemy_id = e.id
LEFT JOIN enemy_skills es ON es.id = r.enemy_skill_id
ORDER BY e.id, es.id;
```

保存前検証:

```sql
SELECT input.enemy_skill_id
FROM input_enemy_skill_ids input
LEFT JOIN enemy_skills es ON es.id = input.enemy_skill_id
WHERE es.id IS NULL;

SELECT input.kind, input.ref_id
FROM input_enemy_drop_table_rows input
LEFT JOIN items i ON input.kind = 'item' AND i.id = input.ref_id
LEFT JOIN grimoire g ON input.kind = 'grimoire' AND g.id = input.ref_id
WHERE (input.kind = 'item' AND i.id IS NULL)
   OR (input.kind = 'grimoire' AND g.id IS NULL);

-- 実装上、POST /api/enemies の ValidateSave は enemySkillIds と
-- 任意の埋め込み dropTable[] を items / grimoire に照合する。
-- drop_table_id -> treasures.id の存在確認は bundle validation 側で実施される。
```

保存:

```sql
MERGE INTO enemies AS dst
USING (SELECT :id, :name, :drop_table_id, :updated_at) AS src(id, name, drop_table_id, updated_at)
ON dst.id = src.id
WHEN MATCHED THEN UPDATE SET name = src.name, drop_table_id = src.drop_table_id, updated_at = src.updated_at
WHEN NOT MATCHED THEN INSERT (id, name, drop_table_id, updated_at) VALUES (src.id, src.name, src.drop_table_id, src.updated_at);

DELETE FROM enemy_skill_refs
WHERE enemy_id = :id;

INSERT INTO enemy_skill_refs (enemy_id, enemy_skill_id)
VALUES (:id, :enemy_skill_id);
```

補足:

```sql
-- UI 初期値では enemies.id = enemies.drop_table_id を前提に入力補助する
SELECT :new_enemy_id AS id, :new_enemy_id AS default_drop_table_id;

-- ただし DB 制約として強制されているわけではなく、
-- 実際には treasures.id を自由入力/参照している
```

## 4. API ごとのデータフロー

画面系と同じ保存先を使い、JSON 入出力にしたものが `/api/*`。

### GET API

```sql
GET /api/items         -> SELECT * FROM items;
GET /api/grimoire      -> SELECT * FROM grimoire;
GET /api/skills        -> SELECT * FROM skills;
GET /api/enemy-skills  -> SELECT * FROM enemy_skills;
GET /api/treasures     -> SELECT * FROM treasures + treasure_loot_pools;
GET /api/enemies       -> SELECT * FROM enemies + enemy_skill_refs;
```

### POST API

```sql
POST /api/items         -> UPSERT items
POST /api/grimoire      -> UPSERT grimoire
POST /api/skills        -> JOIN items で参照確認後 UPSERT skills
POST /api/enemy-skills  -> UPSERT enemy_skills
POST /api/treasures     -> JOIN items/grimoire で参照確認後 UPSERT treasures + treasure_loot_pools
POST /api/enemies       -> JOIN enemy_skills/items/grimoire で検証後 UPSERT enemies + enemy_skill_refs
```

### DELETE API

```sql
DELETE /api/items/:id
DELETE /api/grimoire/:id
DELETE /api/skills/:id
DELETE /api/treasures/:id
DELETE /api/enemies/:id
```

例外:

```sql
-- enemy skill は参照中なら削除不可
DELETE /api/enemy-skills/:id
WHERE NOT EXISTS (
  SELECT 1
  FROM enemy_skill_refs
  WHERE enemy_skill_id = :id
);
```

## 5. 全体検証

`application.ValidateBundle` は全 state をまとめて読み、参照整合性を横断チェックする。

```sql
WITH
  item_ids AS (SELECT id FROM items),
  grimoire_ids AS (SELECT id FROM grimoire),
  enemy_skill_ids AS (SELECT id FROM enemy_skills),
  treasure_ids AS (SELECT id FROM treasures)
SELECT 'skills.item_id missing'
FROM skills s
LEFT JOIN item_ids i ON i.id = s.item_id
WHERE i.id IS NULL

UNION ALL

SELECT 'treasure_loot_pools.ref_id missing(item)'
FROM treasure_loot_pools p
LEFT JOIN item_ids i ON p.kind = 'item' AND i.id = p.ref_id
WHERE p.kind = 'item' AND i.id IS NULL

UNION ALL

SELECT 'treasure_loot_pools.ref_id missing(grimoire)'
FROM treasure_loot_pools p
LEFT JOIN grimoire_ids g ON p.kind = 'grimoire' AND g.id = p.ref_id
WHERE p.kind = 'grimoire' AND g.id IS NULL

UNION ALL

SELECT 'enemy.drop_table_id missing'
FROM enemies e
LEFT JOIN treasure_ids t ON t.id = e.drop_table_id
WHERE t.id IS NULL

UNION ALL

SELECT 'enemy_skill_refs.enemy_skill_id missing'
FROM enemy_skill_refs r
LEFT JOIN enemy_skill_ids es ON es.id = r.enemy_skill_id
WHERE es.id IS NULL;
```

## 6. Export のデータフロー

`/save` または `POST /api/save` で全 state を集約し、export settings を読んで datapack 出力する。

```sql
SELECT *
FROM items;

SELECT *
FROM grimoire;

SELECT *
FROM skills;

SELECT *
FROM enemy_skills;

SELECT *
FROM treasures t
JOIN treasure_loot_pools p ON p.treasure_id = t.id;

SELECT *
FROM enemies e
LEFT JOIN enemy_skill_refs r ON r.enemy_id = e.id;
```

export 時の実質的な join:

```sql
SELECT
  e.id AS enemy_id,
  t.id AS treasure_id,
  p.kind,
  p.ref_id,
  es.id AS enemy_skill_id
FROM enemies e
LEFT JOIN treasures t ON t.id = e.drop_table_id
LEFT JOIN treasure_loot_pools p ON p.treasure_id = t.id
LEFT JOIN enemy_skill_refs r ON r.enemy_id = e.id
LEFT JOIN enemy_skills es ON es.id = r.enemy_skill_id;
```

## 7. 現状の整理ポイント

```sql
-- 1. 主キーは基本的に全 JSON エントリで id(UUID)
SELECT 'items.id / grimoire.id / skills.id / enemy_skills.id / treasures.id / enemies.id';

-- 2. FK 相当は 4 系統
SELECT 'skills.item_id -> items.id';
SELECT 'treasure_loot_pools.ref_id -> items.id | grimoire.id';
SELECT 'enemies.drop_table_id -> treasures.id';
SELECT 'enemy_skill_refs.enemy_skill_id -> enemy_skills.id';

-- 3. enemies と treasures は UI 上かなり近い
SELECT 'defaultEnemyForm では drop_table_id = enemy.id を初期値にしている';

-- 4. ただし保存構造上は enemies と treasures は別 JSON
SELECT 'enemy-state.json と treasure-state.json を id 参照でつないでいる';
```
