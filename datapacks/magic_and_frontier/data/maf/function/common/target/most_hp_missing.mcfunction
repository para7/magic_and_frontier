$tag @e[tag=maf_heal_target] remove maf_heal_target
$execute as @e[type=#maf:friendmob,distance=..$(range)] store result score @s mafTgtHP run data get entity @s Health 1
$execute as @e[type=#maf:friendmob,distance=..$(range)] store result score @s mafTgtMaxHP run attribute @s minecraft:max_health get
$execute as @e[type=#maf:friendmob,distance=..$(range)] run scoreboard players operation @s mafTgtDiff = @s mafTgtMaxHP
$execute as @e[type=#maf:friendmob,distance=..$(range)] run scoreboard players operation @s mafTgtDiff -= @s mafTgtHP
$scoreboard players set #max mafTgtDiff 0
$execute as @e[type=#maf:friendmob,distance=..$(range)] if score @s mafTgtDiff > #max mafTgtDiff run scoreboard players operation #max mafTgtDiff = @s mafTgtDiff
$execute as @e[type=#maf:friendmob,distance=..$(range),sort=nearest] unless entity @e[tag=maf_heal_target] if score @s mafTgtDiff = #max mafTgtDiff run tag @s add maf_heal_target
