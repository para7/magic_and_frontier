
execute as @e[distance=..35,type=item] at @s run particle minecraft:cloud ~ ~ ~ 0.3 0.3 0.3 1 10 normal

execute as @e[distance=..35,type=item] run tp ~ ~ ~

# playsound minecraft:entity.player.levelup player @a ~ ~ ~ 1 2.0

playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5

tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は スウィープ を唱えた！"}]
