effect give @e[distance=..64,type=#p7b:friendmob] minecraft:instant_health 1 1
effect give @e[distance=..64,type=#p7b:friendmob] minecraft:regeneration 10 0
execute at @e[distance=..64,type=#p7b:friendmob] run particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 6
tellraw @a[distance=..64] [{"selector":"@s"},{"text":" は エリアヒール を唱えた！"}]

function maf:magic/exec/00_heal/effect
