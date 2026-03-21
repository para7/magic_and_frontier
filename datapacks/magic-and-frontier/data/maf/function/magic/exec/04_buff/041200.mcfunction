effect give @e[distance=..12,type=#p7b:friendmob] minecraft:strength 60 4
effect give @e[distance=..12,type=#p7b:friendmob] minecraft:resistance 60 2
effect give @e[distance=..12,type=#p7b:friendmob] minecraft:speed 60 3
effect give @e[distance=..12,type=#p7b:friendmob] minecraft:jump_boost 60 1
effect give @e[distance=..12,type=#p7b:friendmob] minecraft:water_breathing 60 0
effect give @e[distance=..12,type=#p7b:friendmob] minecraft:regeneration 60 0
# effect give @e[distance=..12,type=#p7b:friendmob] minecraft:slow_falling 60 0


tellraw @a[distance=..14] [{"selector":"@s"},{"text":" は "},{"nbt":"magictmp.title","storage":"p7:maf"},{"text": "を唱えた！"}]
