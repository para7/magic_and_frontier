execute as @e[type=#p7b:enemymob,distance=..30,nbt={OnGround:1b}] run damage @s 20 minecraft:fall by @p

tellraw @a[distance=..35] [{"selector":"@s"},{"text":" は "},{"nbt":"magictmp.title","storage":"p7:maf"},{"text": "を唱えた！"}]


function maf:magic/exec/01_attack/effect
