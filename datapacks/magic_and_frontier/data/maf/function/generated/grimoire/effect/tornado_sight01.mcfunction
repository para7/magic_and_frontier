function maf:common/sight/eyes_circle_tagged {forward:4.0,radius:3.5,particleCount:80}
execute as @e[type=#maf:enemymob,tag=maf_sight_circle_target] run data modify entity @s Motion set value [0.0d, 1.8d, 0.0d]
execute as @e[type=#maf:enemymob,tag=maf_sight_circle_target] at @s run particle minecraft:cloud ~ ~0.2 ~ 0.75 0.1 0.75 0.01 10 force
execute as @e[type=#maf:enemymob,tag=maf_sight_circle_target] at @s run particle minecraft:campfire_cosy_smoke ~ ~0.5 ~ 0.75 0.1 0.75 0.01 6 force
playsound minecraft:entity.ender_dragon.flap master @a ~ ~ ~ 1.5 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は トルネードサイト を唱えた！"}]
