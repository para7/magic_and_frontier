#ヒーリング
# execute as @e[distance=..14] if entity @s[type=#p7b:undead] if block ~ ~-1 ~ minecraft:blue_ice run 
# effect give @a[distance=..10] minecraft:instant_health 1 1
execute as @e[distance=..14,type=#p7b:enemymob] run data merge entity @s {Motion:[0.0,1.32,0.0]}
execute as @e[distance=..14,type=#p7b:enemymob] at @s run particle minecraft:cloud ~ ~ ~ 0.5 0.5 0.5 0.2 50 normal
playsound minecraft:entity.enderman.teleport master @s ~ ~ ~ 1 1.2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は トルネード を唱えた！"}]
