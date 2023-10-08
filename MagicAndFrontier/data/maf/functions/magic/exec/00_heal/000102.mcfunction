effect give @e[distance=..10,type=#p7b:friendmob] minecraft:absorption 60 4
execute at @e[distance=..10,type=#p7b:friendmob] run particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 6

tellraw @a[distance=..10] [{"selector":"@s"},{"text":" は バリア を唱えた！"}]

function maf:magic/exec/00_heal/effect
