$tag @e[tag=maf_heal_target] remove maf_heal_target
$execute as @e[type=#maf:friendmob,distance=..$(range)] store result score @s mafTgtMaxHP run attribute @s minecraft:max_health get
$scoreboard players set #min mafTgtMaxHP 99999
$execute as @e[type=#maf:friendmob,distance=..$(range)] if score @s mafTgtMaxHP < #min mafTgtMaxHP run scoreboard players operation #min mafTgtMaxHP = @s mafTgtMaxHP
$execute as @e[type=#maf:friendmob,distance=..$(range),sort=nearest] unless entity @e[tag=maf_heal_target] if score @s mafTgtMaxHP = #min mafTgtMaxHP run tag @s add maf_heal_target
