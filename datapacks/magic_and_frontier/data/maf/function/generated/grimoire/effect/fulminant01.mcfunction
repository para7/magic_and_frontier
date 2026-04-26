tag @e[distance=..15,type=#maf:enemymob,tag=maf_target] remove maf_target
execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] run tag @s add maf_target
execute as @e[distance=..15,type=#maf:enemymob,tag=maf_target,limit=1] at @s run execute if block ~ ~ ~ #maf:water as @e[distance=..5.4] run damage @s 192 minecraft:lightning_bolt
execute as @e[distance=..15,type=#maf:enemymob,tag=maf_target,limit=1] at @s run execute unless block ~ ~ ~ #maf:water as @e[distance=..1.8] run damage @s 48 minecraft:lightning_bolt
execute as @e[distance=..15,type=#maf:enemymob,tag=maf_target,limit=1] at @s run damage @s 32 minecraft:lightning_bolt
#execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run summon minecraft:lightning_bolt ~ ~2 ~
#execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run execute if block ~ ~ ~ #maf:water run summon minecraft:lightning_bolt ~2 ~1 ~2
#execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run execute if block ~ ~ ~ #maf:water run summon minecraft:lightning_bolt ~-2 ~1 ~2
#execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run execute if block ~ ~ ~ #maf:water run summon minecraft:lightning_bolt ~2 ~1 ~-2
#execute as @e[distance=..15,type=#maf:enemymob,sort=random,limit=1] at @s run execute if block ~ ~ ~ #maf:water run summon minecraft:lightning_bolt ~-2 ~1 ~-2
tag @e[distance=..15,type=#maf:enemymob,tag=maf_target] remove maf_target
playsound minecraft:entity.evoker.cast_spell player @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell player @a ~ ~ ~ 2 0.5
playsound minecraft:entity.lightning_bolt.impact player @a ~ ~ ~ 1 1.3
playsound minecraft:entity.lightning_bolt.thunder player @a ~ ~ ~ 1 1.3
playsound minecraft:entity.generic.explode player @a ~ ~ ~ 1 1.4
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は フルミナント を唱えた！"}]
