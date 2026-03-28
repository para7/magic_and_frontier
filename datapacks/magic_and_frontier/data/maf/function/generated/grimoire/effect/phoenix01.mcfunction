
# instant_healthをまく
effect give @e[type=#p7b:undead,distance=..70] minecraft:instant_health 1 3
effect give @e[type=#p7b:friendmob,distance=..70] minecraft:instant_health 1 7
effect give @e[type=#p7b:friendmob,distance=..70] minecraft:regeneration 50 1


execute at @e[type=#p7b:undead,distance=..60] run playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 2.0
execute at @e[type=#p7b:friendmob,distance=..60] run playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 2.0

execute at @e[type=#p7b:undead,distance=..30] run particle minecraft:heart ~ ~1 ~ 1.2 1.2 1.2 1 14
execute at @e[type=#p7b:friendmob,distance=..30] run particle minecraft:heart ~ ~1 ~ 1.2 1.2 1.2 1 14

playsound minecraft:entity.evoker.cast_spell player @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell player @a ~ ~ ~ 2 0.5

tellraw @a[distance=..90] [{"selector":"@s"},{"text":" は フェニックス を唱えた！"}]