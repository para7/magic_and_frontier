effect give @e[distance=..6,type=#p7b:friendmob] minecraft:instant_health 1 0
effect give @e[distance=..6,type=#p7b:friendmob] minecraft:regeneration 10 0
execute at @e[distance=..6,type=#p7b:friendmob] run particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 6
tellraw @a[distance=..6] [{"selector":"@s"},{"text":" は ヒーリング を唱えた！"}]

function maf:magic/exec/00_heal/effect
