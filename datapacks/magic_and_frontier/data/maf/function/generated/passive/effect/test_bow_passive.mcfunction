execute unless score @s mafBowUsed matches 1.. run return 0
execute store result storage maf:tmp bow_player_id int 1 run scoreboard players get @s mafPlayerID
execute as @e[type=arrow,distance=..5,nbt=!{inGround:1b},sort=nearest,limit=1] run function maf:magic/passive/tag_passive_arrow {passive_id:"test_bow_passive",life:1100}
