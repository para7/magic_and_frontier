$tag @e[tag=maf_heal_target] remove maf_heal_target
$execute as @e[type=#maf:friendmob,distance=..$(range)] store result score @s mafTgtArmor run attribute @s minecraft:armor_toughness get
$scoreboard players set #max mafTgtArmor 0
$execute as @e[type=#maf:friendmob,distance=..$(range)] if score @s mafTgtArmor > #max mafTgtArmor run scoreboard players operation #max mafTgtArmor = @s mafTgtArmor
$execute as @e[type=#maf:friendmob,distance=..$(range),sort=nearest] unless entity @e[tag=maf_heal_target] if score @s mafTgtArmor = #max mafTgtArmor run tag @s add maf_heal_target
