effect give @s minecraft:levitation 1 9
effect give @s minecraft:slow_falling 3 0

execute at @e[distance=..4,type=#p7b:friendmob] run particle minecraft:glow ~ ~1 ~ 0.5 0.5 0.5 1 8
tellraw @a[distance=..4] [{"selector":"@s"},{"text":" は レビテート を唱えた！"}]

function maf:magic/exec/02_live/effect
