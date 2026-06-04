execute as @e[type=#maf:friendmob,distance=..10] run effect give @s minecraft:health_boost 300 1
execute at @e[type=#maf:friendmob,distance=..10] run particle minecraft:heart ~ ~1 ~ 0.3 0.3 0.3 1 6
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は ハートブースト を唱えた！"}]
