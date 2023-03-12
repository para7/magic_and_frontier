#ヒーリング


effect give @e[distance=..6,type=#para7_utils:friendmob] minecraft:instant_health 1 0
effect give @e[distance=..6,type=#para7_utils:friendmob] minecraft:regeneration 10 0

execute at @e[distance=..6,type=#para7_utils:friendmob] run particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 6

# particle minecraft:heart ~ ~ ~ 1 1 1 1 5

playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 2.0

playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5

tellraw @a[distance=..6] [{"selector":"@s"},{"text":" は ヒーリング を唱えた！"}]

scoreboard players remove @s p7_MP 15