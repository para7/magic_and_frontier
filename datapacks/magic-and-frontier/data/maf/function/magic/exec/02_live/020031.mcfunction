execute as @e[type=item,distance=..30] run teleport ~ ~ ~

# execute at @e[distance=..4,type=#p7b:friendmob] run particle minecraft:glow ~ ~1 ~ 0.5 0.5 0.5 1 8
tellraw @a[distance=..4] [{"selector":"@s"},{"text":" は "},{"nbt":"magictmp.title","storage":"p7:maf"},{"text": "を唱えた！"}]

function maf:magic/exec/02_live/effect
