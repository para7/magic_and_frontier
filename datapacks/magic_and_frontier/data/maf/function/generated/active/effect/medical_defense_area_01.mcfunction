execute as @e[type=#maf:friendmob,distance=..10] run effect give @s minecraft:resistance 120 3
execute at @e[type=#maf:friendmob,distance=..10] run particle minecraft:enchanted_hit ~ ~1 ~ 0.3 0.5 0.3 0.05 12
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は 医術防御 を唱えた！"}]
