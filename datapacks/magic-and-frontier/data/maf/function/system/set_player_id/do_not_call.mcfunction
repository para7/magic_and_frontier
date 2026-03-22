tellraw @a "プレイヤーID付与処理(MPバー)"

# IDをインクリメントして保存する
execute store result score @s mafPlayerID run data get storage p7:maf playerSequence

scoreboard players add @s mafPlayerID 1

execute store result storage p7:maf playerSequence byte 1 run scoreboard players get @s mafPlayerID


# execute store result storage p7:maf playerSequence int 1 run scoreboard players get @s mafPlayerID