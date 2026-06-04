execute as @e[type=#maf:friendmob,distance=0.1..10] run effect clear @s minecraft:poison
execute as @e[type=#maf:friendmob,distance=0.1..10] run effect clear @s minecraft:slowness
execute as @e[type=#maf:friendmob,distance=0.1..10] run effect clear @s minecraft:mining_fatigue
execute as @e[type=#maf:friendmob,distance=0.1..10] run effect clear @s minecraft:blindness
execute as @e[type=#maf:friendmob,distance=0.1..10] run effect clear @s minecraft:wither
execute as @e[type=#maf:friendmob,distance=0.1..10] run effect clear @s minecraft:hunger
execute as @e[type=#maf:friendmob,distance=0.1..10] run effect clear @s minecraft:levitation
execute at @e[type=#maf:friendmob,distance=0.1..10] run particle minecraft:happy_villager ~ ~1 ~ 0.3 0.3 0.3 1 6
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.8
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は リフレッシュ を唱えた！"}]
