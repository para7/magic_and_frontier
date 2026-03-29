#アブソリュートゼロ
fill ~-15 ~-3 ~-15 ~15 ~1 ~15 minecraft:frosted_ice replace water

# playsound minecraft:entity.blaze.shoot master @a ~ ~ ~ 1 1
# playsound minecraft:entity.generic.explode master @a ~ ~ ~

# function para7_utils:set_undead_tag


# あとでこれを入れる
# unless entity @s[type=minecraft:player] 

# execute as @a[distance=..18] at @s unless entity @s[type=minecraft:player] if entity @s[type=!#p7b:undead] if block ~ ~-1 ~ minecraft:frosted_ice run effect give @s minecraft:instant_damage 1 1

execute as @e[distance=..18] at @s unless entity @s[type=minecraft:player] if entity @s[type=!#p7b:undead] if block ~ ~-1 ~ minecraft:frosted_ice run effect give @s minecraft:instant_damage 1 1
execute as @e[distance=..18] at @s unless entity @s[type=minecraft:player] if entity @s[type=!#p7b:undead] if block ~ ~ ~ minecraft:frosted_ice run effect give @s minecraft:instant_damage 1 3
execute as @e[distance=..18] at @s unless entity @s[type=minecraft:player] if entity @s[type=!#p7b:undead] if block ~ ~-1 ~ minecraft:ice run effect give @s minecraft:instant_damage 1 3
execute as @e[distance=..18] at @s unless entity @s[type=minecraft:player] if entity @s[type=!#p7b:undead] if block ~ ~-1 ~ minecraft:packed_ice run effect give @s minecraft:instant_damage 1 5
execute as @e[distance=..18] at @s unless entity @s[type=minecraft:player] if entity @s[type=!#p7b:undead] if block ~ ~-1 ~ minecraft:blue_ice run effect give @s minecraft:instant_damage 1 8


execute as @e[distance=..18] at @s unless entity @s[type=minecraft:player] if entity @s[type=#p7b:undead] if block ~ ~-1 ~ minecraft:frosted_ice run effect give @s minecraft:instant_health 1 1
execute as @e[distance=..18] at @s unless entity @s[type=minecraft:player] if entity @s[type=#p7b:undead] if block ~ ~ ~ minecraft:frosted_ice run effect give @s minecraft:instant_health 1 3
execute as @e[distance=..18] at @s unless entity @s[type=minecraft:player] if entity @s[type=#p7b:undead] if block ~ ~-1 ~ minecraft:ice run effect give @s minecraft:instant_health 1 3
execute as @e[distance=..18] at @s unless entity @s[type=minecraft:player] if entity @s[type=#p7b:undead] if block ~ ~-1 ~ minecraft:packed_ice run effect give @s minecraft:instant_health 1 5
execute as @e[distance=..18] at @s unless entity @s[type=minecraft:player] if entity @s[type=#p7b:undead] if block ~ ~-1 ~ minecraft:blue_ice run effect give @s minecraft:instant_health 1 8

execute as @a[distance=..18] at @s if block ~ ~-1 ~ minecraft:frosted_ice run effect give @s minecraft:instant_damage 1 0
execute as @a[distance=..18] at @s if block ~ ~ ~ minecraft:frosted_ice run effect give @s minecraft:instant_damage 1 1
execute as @a[distance=..18] at @s if block ~ ~-1 ~ minecraft:ice run effect give @s minecraft:instant_damage 1 0
execute as @a[distance=..18] at @s if block ~ ~-1 ~ minecraft:packed_ice run effect give @s minecraft:instant_damage 1 0
execute as @a[distance=..18] at @s if block ~ ~-1 ~ minecraft:blue_ice run effect give @s minecraft:instant_damage 1 0




tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は アブソリュート・ゼロ を唱えた！"}]

playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5
playsound minecraft:block.glass.break master @a ~ ~ ~ 3 0.6