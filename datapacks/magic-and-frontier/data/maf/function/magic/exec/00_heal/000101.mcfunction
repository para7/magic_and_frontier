effect give @e[distance=..10,type=#p7b:friendmob] minecraft:absorption 14 0
execute at @e[distance=..10,type=#p7b:friendmob] run particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 6

tellraw @a[distance=..10] [{"selector":"@s"},{"text":" は "},{"nbt":"magictmp.title","storage":"p7:maf"},{"text": "を唱えた！"}]


function maf:magic/exec/00_heal/effect
