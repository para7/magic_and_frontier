tellraw @p [{"selector":"@s"},{"text":" は "},{"nbt":"magictmp.title","storage":"p7:maf"},{"text": "を唱えた！"}]

# 詠唱時のストレージ巻き上げ処理を利用してアナライズを実行する
# バグ起きるかも. 要注意
scoreboard players set @s mafCastID 0
# 対象外アイテムの場合は0がセットされる
execute store result score @s mafCastID run data get entity @s Inventory[{Slot:9b}].tag.magicID

execute if entity @s[scores={mafCastID=0}] run tellraw @s [{"text":"インベントリの左上に対象アイテムをセットしてください"}]
execute unless entity @s[scores={mafCastID=0}] run function maf:magic/exec/02_live/029001_2


# modifires 置き換え
function maf:magic/exec/generated/analyze

# 終了処理
scoreboard players set @s mafCastID 0
function maf:magic/exec/02_live/effect
