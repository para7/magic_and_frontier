execute as @e[distance=..20,type=#p7b:enemymob] run damage @s 6 minecraft:in_fire by @p

execute as @e[distance=..20,type=#p7b:enemymob] run data merge entity @s {Fire: 200}

# execute as @e[distance=..10] run particle minecraft:heart ~ ~1 ~ 0.5 0.5 0.5 1 6

tellraw @a[distance=..25] [{"selector":"@s"},{"text":" は ファイアストーム を唱えた！"}]

function maf:magic/exec/01_attack/effect
