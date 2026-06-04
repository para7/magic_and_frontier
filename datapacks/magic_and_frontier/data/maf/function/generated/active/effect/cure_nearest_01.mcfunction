execute as @e[type=#maf:friendmob,distance=0.1..10,sort=nearest,limit=1] run effect give @s minecraft:instant_health 1 0
execute at @e[type=#maf:friendmob,distance=0.1..10,sort=nearest,limit=1] run particle minecraft:heart ~ ~1 ~ 0.3 0.3 0.3 1 4
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 2.0
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は キュア を唱えた！"}]
