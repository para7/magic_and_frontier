effect give @e[distance=..6,type=#p7b:friendmob] minecraft:haste 40 3
# effect give @e[distance=..6,type=#p7b:friendmob] minecraft:regeneration 10 0
execute at @e[distance=..6,type=#p7b:friendmob] run particle minecraft:glow ~ ~1 ~ 0.5 0.5 0.5 1 8
tellraw @a[distance=..6] [{"selector":"@s"},{"text":" は ヘイスト を唱えた！"}]

# function maf:magic/exec/02_live/effect
