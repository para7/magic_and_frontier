tag @e[tag=maf_target] remove maf_target
execute store result score @s tmp run execute as @e[distance=..15,sort=nearest,limit=1,type=#maf:enemymob] at @s run tag @e[distance=..5,type=#maf:enemymob] add maf_target
tellraw @p [{"text":"[DEBUG] chainbolt hit count: "},{"score":{"name":"@s","objective":"tmp"}}]
execute if score @s tmp matches 1.. as @e[type=#maf:enemymob,tag=maf_target] run damage @s 10 minecraft:lightning_bolt
execute if score @s tmp matches 2.. as @e[type=#maf:enemymob,tag=maf_target] run damage @s 10 minecraft:lightning_bolt
execute if score @s tmp matches 3.. as @e[type=#maf:enemymob,tag=maf_target] run damage @s 10 minecraft:lightning_bolt
execute if score @s tmp matches 4.. as @e[type=#maf:enemymob,tag=maf_target] run damage @s 10 minecraft:lightning_bolt
execute if score @s tmp matches 5.. as @e[type=#maf:enemymob,tag=maf_target] run damage @s 10 minecraft:lightning_bolt
execute if score @s tmp matches 6.. as @e[type=#maf:enemymob,tag=maf_target] run damage @s 10 minecraft:lightning_bolt
execute if score @s tmp matches 7.. as @e[type=#maf:enemymob,tag=maf_target] run damage @s 10 minecraft:lightning_bolt
execute if score @s tmp matches 8.. as @e[type=#maf:enemymob,tag=maf_target] run damage @s 10 minecraft:lightning_bolt
execute if score @s tmp matches 9.. as @e[type=#maf:enemymob,tag=maf_target] run damage @s 15 minecraft:lightning_bolt
execute if score @s tmp matches 10.. as @e[type=#maf:enemymob,tag=maf_target] run damage @s 15 minecraft:lightning_bolt
tag @e[tag=maf_target] remove maf_target
playsound minecraft:entity.evoker.cast_spell player @a ~ ~ ~ 1.2 1.3
execute if score @s tmp matches 1..2 run playsound minecraft:entity.lightning_bolt.impact player @a ~ ~ ~ 0.6 1.8
execute if score @s tmp matches 3.. run playsound minecraft:entity.lightning_bolt.impact player @a ~ ~ ~ 1.0 1.3
#execute if score @s tmp matches 5..8 run playsound minecraft:entity.generic.explode player @a ~ ~ ~ 1.0 1.5
execute if score @s tmp matches 7.. run playsound minecraft:entity.lightning_bolt.thunder player @a ~ ~ ~ 1.2 1.2
#execute if score @s tmp matches 9.. run playsound minecraft:entity.generic.explode player @a ~ ~ ~ 1.8 0.8
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は チェインボルト を唱えた！"}]
