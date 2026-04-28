execute as @e[type=#maf:enemymob,distance=..8] run data modify entity @s Motion set value [0.0d, 1.8d, 0.0d]
execute at @s run particle minecraft:cloud ~ ~0.2 ~ 4.0 0.1 4.0 0.01 250 force
execute at @s run particle minecraft:campfire_cosy_smoke ~ ~0.5 ~ 3.0 0.1 3.0 0.01 120 force
playsound minecraft:entity.ender_dragon.flap master @a ~ ~ ~ 1.5 1.4
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は トルネード を唱えた！"}]
