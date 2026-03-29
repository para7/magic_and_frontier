
playsound minecraft:entity.ender_dragon.ambient player @a ~ ~ ~ 0.8 1.5
playsound minecraft:entity.evoker.cast_spell player @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell player @a ~ ~ ~ 2 0.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は テレポート を唱えた！"}]

teleport @e[distance=..5] ^ ^ ^200

playsound minecraft:entity.ender_dragon.ambient player @a ^ ^ ^200 0.8 1.5
playsound minecraft:entity.evoker.cast_spell player @a ^ ^ ^200 2 2
playsound minecraft:entity.evoker.cast_spell player @a ^ ^ ^200 2 0.5
# tellraw @a[distance=1..30] [{"selector":"@s"},{"text":" は テレポート を唱えた！"}]

# xp add @s -1 levels