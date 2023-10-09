effect give @e[distance=..6,type=#p7b:friendmob] minecraft:strength 40 3
# effect give @e[distance=..6,type=#p7b:friendmob] minecraft:regeneration 10 0
execute at @e[distance=..6,type=#p7b:friendmob] run particle minecraft:glow ~ ~1 ~ 0.5 0.5 0.5 1 8
tellraw @a[distance=..6] [{"selector":"@s"},{"text":" は "},{"nbt":"magictmp.title","storage":"p7:maf"},{"text": "を唱えた！"}]
