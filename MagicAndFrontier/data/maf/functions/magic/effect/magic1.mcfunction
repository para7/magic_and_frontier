#ヒーリング


# effect give @e[distance=..6] minecraft:instant_health 1 0
# effect give @e[distance=..6] minecraft:regeneration 10 0
# execute at @e[distance=..6] run particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 6

effect give @e[distance=..6,type=#p7b:friendmob] minecraft:instant_health 1 0
effect give @e[distance=..6,type=#p7b:friendmob] minecraft:regeneration 10 0
execute at @e[distance=..6,type=#p7b:friendmob] run particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 6

particle minecraft:heart ~ ~ ~ 1 1 1 1 5

playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 2.0

playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5

tellraw @a[distance=..6] [{"selector":"@s"},{"text":" は ヒーリング を唱えた！"}]