execute at @s anchored eyes rotated as @s run summon minecraft:area_effect_cloud ^ ^ ^4 {Tags:["maf_tornado_sight_center"],Duration:8,Radius:0.0f,WaitTime:0,ReapplicationDelay:0}
execute as @e[type=#maf:enemymob,distance=..20] at @s if entity @e[type=minecraft:area_effect_cloud,tag=maf_tornado_sight_center,distance=..3.5,limit=1] run data modify entity @s Motion set value [0.0d, 1.8d, 0.0d]
execute at @e[type=minecraft:area_effect_cloud,tag=maf_tornado_sight_center,distance=..20,sort=nearest,limit=1] run particle minecraft:cloud ~ ~0.2 ~ 1.8 0.1 1.8 0.01 100 force
execute at @e[type=minecraft:area_effect_cloud,tag=maf_tornado_sight_center,distance=..20,sort=nearest,limit=1] run particle minecraft:campfire_cosy_smoke ~ ~0.5 ~ 1.8 0.1 1.8 0.01 60 force
kill @e[type=minecraft:area_effect_cloud,tag=maf_tornado_sight_center,distance=..20]
playsound minecraft:entity.ender_dragon.flap master @a ~ ~ ~ 1.5 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は トルネードサイト を唱えた！"}]
