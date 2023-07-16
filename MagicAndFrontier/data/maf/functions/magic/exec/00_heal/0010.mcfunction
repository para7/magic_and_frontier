effect give @e[distance=..6,type=#p7b:friendmob] minecraft:absorption 20 0
execute at @e[distance=..6,type=#p7b:friendmob] run particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 6

tellraw @a[distance=..6] [{"selector":"@s"},{"text":" は バリア を唱えた！"}]

function maf:magic/exec/00_heal/effect
