execute unless score @s mafBowUsed matches 1.. run return 0
execute store result storage maf:tmp bow_player_id int 1 run scoreboard players get @s mafPlayerID
execute as @e[type=arrow,distance=..2,nbt=!{inGround:1b},sort=nearest,limit=1] run function maf:bow/tag_bow_arrow {bow_id:"test_full",life:1100}
execute as @e[type=arrow,distance=..2,tag=maf_bow_arrow,sort=nearest,limit=1] run function maf:bow/prepare_hit_arrow
execute as @e[type=arrow,distance=..2,tag=maf_bow_arrow,sort=nearest,limit=1] run tag @s add flying
execute as @e[type=arrow,distance=..2,tag=maf_bow_arrow,sort=nearest,limit=1] run tag @s add ground
say fired
