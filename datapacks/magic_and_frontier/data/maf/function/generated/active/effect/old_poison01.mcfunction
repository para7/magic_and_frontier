#テンペスト
# function 
execute as @e[distance=..8,type=#maf:enemymob,sort=nearest,limit=1] at @s run effect give @e[type=#maf:undead,distance=..2.9] minecraft:instant_health 1 0
execute as @e[distance=..8,type=#maf:enemymob,sort=nearest,limit=1] at @s run effect give @e[type=!#maf:undead,distance=..2.9] minecraft:instant_damage 1 1
execute as @e[distance=..8,type=#maf:enemymob,sort=nearest,limit=1] at @s run effect give @e[distance=..2.9] minecraft:poison 13 3
execute as @e[distance=..8,type=#maf:enemymob,sort=nearest,limit=1] at @s run effect give @e[distance=..2.9] minecraft:weakness 30 3
execute as @e[distance=..8,type=#maf:enemymob,sort=nearest,limit=1] at @s run effect clear @e[distance=..2.9] minecraft:resistance
# execute as @e[distance=..20,tag=MOB_NotFriend,sort=nearest,limit=5] at @s run effect give @s glowing
playsound minecraft:entity.silverfish.hurt player @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5
# playsound minecraft:entity.lightning_bolt.impact weather @a ~ ~ ~ 5 1
# playsound minecraft:entity.lightning_bolt.thunder weather @a ~ ~ ~ 5 1
tellraw @a[distance=..18] [{"selector":"@s"},{"text":" は ポイズン を唱えた！"}]
