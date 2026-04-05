---
name: ohmydat
description: Use this skill when working with the oh_my_dat datapack or when you need to read, write, or delete per-entity private storage in Minecraft commands.
---

## 使い方/How To Use

### 基本 / Basics

```mcfunction
個別ストレージを使いたいエンティティで次を実行するだけ！ / Run the following command as the entity you want to use private storage for.
function #oh_my_dat:please

アクセス / Access
data get storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].DataName

書き換え / Modify
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].DataName set value DataValue

消去 / Delete
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].DataName
```

### 既知のストレージにIDでアクセスしたい場合 / If you want to use known storage by storage ID

```mcfunction
scoreboard players set _ OhMyDatID <ID>
function #oh_its_dat:please
```