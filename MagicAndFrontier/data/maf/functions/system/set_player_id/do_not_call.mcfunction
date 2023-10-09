tellraw @a "プレイヤーID付与処理(MPバー)"

# IDをインクリメントして保存する
execute store result score @s p7_playerID run data get storage p7:maf playerSequence

scoreboard players add @s p7_playerID 1

execute store result storage p7:maf playerSequence byte 1 run scoreboard players get @s p7_playerID


# execute store result storage p7:maf playerSequence int 1 run scoreboard players get @s p7_playerID