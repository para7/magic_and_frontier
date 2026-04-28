execute at @s anchored eyes rotated as @s run summon minecraft:area_effect_cloud ^ ^ ^4 {Tags:["maf_poison_sight_center"],Duration:8,Radius:0.0f,WaitTime:0,ReapplicationDelay:0}
execute as @e[type=#maf:enemymob,distance=..20] at @s if entity @e[type=minecraft:area_effect_cloud,tag=maf_poison_sight_center,distance=..3.5,limit=1] run effect give @s minecraft:poison 10 10
execute at @e[type=minecraft:area_effect_cloud,tag=maf_poison_sight_center,distance=..20,sort=nearest,limit=1] run particle minecraft:witch ~ ~0.2 ~ 1.8 0.1 1.8 0.01 80 force
execute at @e[type=minecraft:area_effect_cloud,tag=maf_poison_sight_center,distance=..20,sort=nearest,limit=1] run particle minecraft:spore_blossom_air ~ ~0.2 ~ 1.8 0.2 1.8 0.01 32 force
kill @e[type=minecraft:area_effect_cloud,tag=maf_poison_sight_center,distance=..20]
playsound minecraft:entity.witch.throw master @a ~ ~ ~ 1.5 1.1
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は ポイズンサイト を唱えた！"}]
