execute as @e[type=#maf:undead,distance=..8] run damage @s 20 minecraft:magic
execute at @s run particle minecraft:soul_fire_flame ~ ~0.2 ~ 4.0 0.2 4.0 0.01 200 force
execute as @e[type=#maf:undead,distance=..8] at @s run particle minecraft:soul ~ ~0.5 ~ 0.5 0.5 0.5 0.01 5 force
playsound minecraft:block.bell.resonate master @a ~ ~ ~ 2.0 0.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は レクイエム を唱えた！"}]
