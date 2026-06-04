$tag @e[tag=maf_heal_target] remove maf_heal_target
$execute as @e[type=#maf:friendmob,distance=..$(range)] store result score @s mafTgtHP run data get entity @s Health 1
$scoreboard players set #min mafTgtHP 99999
$execute as @e[type=#maf:friendmob,distance=..$(range)] if score @s mafTgtHP < #min mafTgtHP run scoreboard players operation #min mafTgtHP = @s mafTgtHP
$execute as @e[type=#maf:friendmob,distance=..$(range),sort=nearest] unless entity @e[tag=maf_heal_target] if score @s mafTgtHP = #min mafTgtHP run tag @s add maf_heal_target
