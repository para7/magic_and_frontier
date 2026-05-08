effect give @e[type=#maf:enemymob,distance=..8] minecraft:poison 10 2
execute at @s run particle minecraft:witch ~ ~0.2 ~ 4.0 0.2 4.0 0.01 200 force
execute at @s run particle minecraft:spore_blossom_air ~ ~0.2 ~ 4.0 0.3 4.0 0.01 200 force
playsound minecraft:entity.witch.throw master @a ~ ~ ~ 1.5 0.8
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は ポイズン を唱えた！"}]
