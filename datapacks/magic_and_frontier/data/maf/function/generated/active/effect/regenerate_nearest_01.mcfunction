execute as @e[type=#maf:friendmob,distance=0.1..10,sort=nearest,limit=1] run effect give @s minecraft:regeneration 30 1
execute at @e[type=#maf:friendmob,distance=0.1..10,sort=nearest,limit=1] run particle minecraft:heart ~ ~1 ~ 0.3 0.3 0.3 0.5 6
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 2.0
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は リジェネレート を唱えた！"}]
