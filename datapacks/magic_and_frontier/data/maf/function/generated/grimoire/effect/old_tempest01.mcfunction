#テンペスト
# function 
execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run effect give @s minecraft:instant_health 1 4
execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run effect give @s minecraft:instant_damage 1 4
execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run execute if block ~ ~ ~ #maf:water run effect give @e[type=#maf:undead,distance=..8.4] minecraft:instant_health 1 5
execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run execute if block ~ ~ ~ #maf:water run effect give @e[type=!#maf:undead,distance=..8.4] minecraft:instant_damage 1 5
execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run execute unless block ~ ~ ~ #maf:water run effect give @e[type=#maf:undead,distance=..1.8] minecraft:instant_health 1 3
execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run execute unless block ~ ~ ~ #maf:water run effect give @e[type=!#maf:undead,distance=..1.8] minecraft:instant_damage 1 3
execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run summon minecraft:lightning_bolt ~ ~2 ~
execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run execute if block ~ ~ ~ #maf:water run summon minecraft:lightning_bolt ~2 ~1 ~2
execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run execute if block ~ ~ ~ #maf:water run summon minecraft:lightning_bolt ~-2 ~1 ~2
execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run execute if block ~ ~ ~ #maf:water run summon minecraft:lightning_bolt ~2 ~1 ~-2
execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run execute if block ~ ~ ~ #maf:water run summon minecraft:lightning_bolt ~-2 ~1 ~-2
# execute as @e[distance=..20,tag=MOB_NotFriend,sort=nearest,limit=5] at @s run effect give @s glowing
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5
# playsound minecraft:entity.lightning_bolt.impact weather @a ~ ~ ~ 5 1
# playsound minecraft:entity.lightning_bolt.thunder weather @a ~ ~ ~ 5 1
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は テンペスト を唱えた！"}]
