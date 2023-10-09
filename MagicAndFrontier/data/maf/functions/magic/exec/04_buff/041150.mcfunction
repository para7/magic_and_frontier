effect give @e[distance=..12,type=#p7b:friendmob] minecraft:strength 60 2
effect give @e[distance=..12,type=#p7b:friendmob] minecraft:speed 60 1
effect give @e[distance=..12,type=#p7b:friendmob] minecraft:jump_boost 60 1


tellraw @a[distance=..12] [{"selector":"@s"},{"text":" は "},{"nbt":"magictmp.title","storage":"p7:maf"},{"text": "を唱えた！"}]
