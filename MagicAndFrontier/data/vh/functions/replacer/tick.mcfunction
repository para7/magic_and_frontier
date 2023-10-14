function p7b:generate_rand

execute if entity @s[type=minecraft:zombie] if score rand p7_Rand1 matches 0..1000 run function vh:summon/bookmaster