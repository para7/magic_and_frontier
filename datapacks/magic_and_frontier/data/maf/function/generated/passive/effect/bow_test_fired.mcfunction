execute unless score @s mafBowUsed matches 1.. run return 0
execute store result storage maf:tmp bow_player_id int 1 run scoreboard players get @s mafPlayerID
execute as @e[type=arrow,distance=..2,nbt=!{inGround:1b},sort=nearest,limit=1] run function maf:magic/bow/tag_bow_arrow {bow_id:"test_fired",life:1100}
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 1 1.5 1
