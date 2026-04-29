function maf:common/sight/eyes_circle_tagged {forward:4.0,radius:3.5,particleCount:80}
effect give @e[type=#maf:enemymob,tag=maf_sight_circle_target] minecraft:poison 10 10
execute as @e[type=#maf:enemymob,tag=maf_sight_circle_target] at @s run particle minecraft:spore_blossom_air ~ ~0.2 ~ 0.75 0.2 0.75 0.01 8 force
playsound minecraft:entity.witch.throw master @a ~ ~ ~ 1.5 1.1
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は ポイズンサイト を唱えた！"}]
