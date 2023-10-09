effect give @s minecraft:strength 120 8
effect give @s minecraft:resistance 120 3
effect give @s minecraft:speed 120 1
effect give @s minecraft:haste 120 2
effect give @s minecraft:jump_boost 120 1
effect give @s minecraft:water_breathing 120 0
effect give @s minecraft:wither 10 0
effect give @s minecraft:saturation 1 5

scoreboard players remove @s p7_soul 50

tellraw @a[distance=..12] [{"selector":"@s"},{"text":" は "},{"nbt":"magictmp.title","storage":"p7:maf"},{"text": "を唱えた！"}]
