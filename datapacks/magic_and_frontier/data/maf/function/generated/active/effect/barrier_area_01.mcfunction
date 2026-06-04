execute as @e[type=#maf:friendmob,distance=..10] run effect give @s minecraft:absorption 210 2
execute at @e[type=#maf:friendmob,distance=..10] run particle minecraft:end_rod ~ ~1 ~ 0.3 0.5 0.3 0.05 8
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は バリア を唱えた！"}]
