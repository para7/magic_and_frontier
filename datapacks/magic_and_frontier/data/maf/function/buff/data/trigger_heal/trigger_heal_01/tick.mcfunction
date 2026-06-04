execute store result score @s mafTgtHP run data get entity @s Health 2
execute store result score @s mafTgtMaxHP run attribute @s minecraft:max_health get
execute if score @s mafTgtHP <= @s mafTgtMaxHP run effect give @s minecraft:instant_health 1 1
execute if score @s mafTgtHP <= @s mafTgtMaxHP at @s run particle minecraft:heart ~ ~1.0 ~ 0.3 0.5 0.3 1 8 force
execute if score @s mafTgtHP <= @s mafTgtMaxHP at @s run playsound minecraft:entity.player.levelup master @a[distance=..24] ~ ~ ~ 1.5 2.0
execute if score @s mafTgtHP <= @s mafTgtMaxHP run data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current.tick set value 1
