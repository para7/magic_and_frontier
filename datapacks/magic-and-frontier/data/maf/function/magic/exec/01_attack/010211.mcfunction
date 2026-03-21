execute as @e[distance=..50,type=#p7b:enemymob] run damage @s 20 minecraft:freeze by @p

effect give @e[distance=..50,type=#p7b:enemymob] minecraft:slowness 40 10

tellraw @a[distance=..25] [{"selector":"@s"},{"text":" は "},{"nbt":"magictmp.title","storage":"p7:maf"},{"text": "を唱えた！"}]


function maf:magic/exec/01_attack/effect
