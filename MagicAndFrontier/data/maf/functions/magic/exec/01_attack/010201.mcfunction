execute as @e[distance=..20,type=#p7b:enemymob] run damage @s 15 minecraft:freeze by @p

effect give @e[distance=..20,type=#p7b:enemymob] minecraft:slowness 30 2

tellraw @a[distance=..25] [{"selector":"@s"},{"text":" は ブリザード を唱えた！"}]

function maf:magic/exec/01_attack/effect
