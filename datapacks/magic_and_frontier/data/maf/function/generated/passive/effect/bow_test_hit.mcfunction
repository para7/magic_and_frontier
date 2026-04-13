execute unless score @s mafBowUsed matches 1.. unless score @s mafCrossbowUsed matches 1.. run return 0
execute store result storage maf:tmp bow_player_id int 1 run scoreboard players get @s mafPlayerID
execute as @e[type=arrow,distance=..2,nbt=!{inGround:1b},sort=nearest,limit=1] run function maf:bow/tag_bow_arrow {bow_id:"test_hit",life:1100}
execute if data entity @s SelectedItem{id:"minecraft:crossbow"} if data entity @s SelectedItem.components."minecraft:enchantments"."minecraft:multishot" as @e[type=arrow,distance=..2,nbt=!{inGround:1b},sort=nearest,limit=3] run function maf:bow/tag_bow_arrow {bow_id:"test_hit",life:1100}
execute as @e[type=arrow,distance=..2,tag=maf_bow_arrow_new] run function maf:bow/prepare_hit_arrow
tag @e[type=arrow,distance=..2,tag=maf_bow_arrow_new] remove maf_bow_arrow_new
