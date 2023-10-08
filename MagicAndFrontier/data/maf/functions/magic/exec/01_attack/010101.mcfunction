execute as @e[type=#p7b:enemymob,distance=..15,nbt={OnGround:1b}] run damage @s 12 minecraft:fall by @p

tellraw @a[distance=..20] [{"selector":"@s"},{"text":" は グランドシェイク を唱えた！"}]

function maf:magic/exec/01_attack/effect

# 前世の魔法↓
# #テンペスト

# # function 

# # execute as @e[distance=..20,type=#para7_utils:enemymob,sort=random,limit=1] at @s run function p_magic:exe/element/magic3_effect

# execute as @e[distance=..20,nbt={OnGround:1b},type=#para7_utils:undead] run effect give @s minecraft:instant_health 1 1
# execute as @e[distance=..20,nbt={OnGround:1b},type=!#para7_utils:undead] run effect give @s minecraft:instant_damage 1 1
# execute as @e[distance=..20,nbt={OnGround:1b},type=!#para7_utils:undead] at @s run particle minecraft:cloud ~ ~ ~ 0.5 0.8 0.5 0.2 50 normal

# playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
# playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5
# playsound minecraft:entity.blaze.shoot master @a ~ ~ ~ 1 1
# playsound minecraft:entity.generic.explode master @a ~ ~ ~

# tellraw @a[distance=..20] [{"selector":"@s"},{"text":" は アースシェイク を唱えた！"}]

# scoreboard players add @s PMag_Enable 1
