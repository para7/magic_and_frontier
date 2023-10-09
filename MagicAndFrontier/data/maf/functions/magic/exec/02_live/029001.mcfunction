tellraw @a[distance=..20] [{"selector":"@s"},{"text":" は "},{"nbt":"magictmp.title","storage":"p7:maf"},{"text": "を唱えた！"}]

# 詠唱時のストレージ巻き上げ処理を利用してアナライズを実行する
# バグ起きるかも. 要注意
execute store result score @s p7_magicID run data get entity @s Inventory[{Slot:9b}].tag.magicID
function maf:magic/exec/selectdb
scoreboard players set @s p7_magicID 0

tellraw @s [{"nbt":"magictmp.title","storage":"p7:maf"}]
tellraw @s [{"text":"詠唱tick: "},{"nbt":"magictmp.cast","storage":"p7:maf"}]
tellraw @s [{"text":"消費MP : "},{"nbt":"magictmp.cost","storage":"p7:maf"}]
tellraw @s [{"text":"効果: "},{"nbt":"magictmp.description","storage":"p7:maf"}]

# modifires 置き換えのサンプル
# /item modify entity @s inventory.0 maf:test