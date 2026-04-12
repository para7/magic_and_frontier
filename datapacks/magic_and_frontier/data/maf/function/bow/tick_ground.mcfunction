# TODO: この @e をどうにか軽量化できないだろうか？
execute as @e[type=arrow,tag=maf_bow_arrow,tag=flying,nbt={inGround:1b}] run tag @s remove flying
execute as @e[type=arrow,tag=maf_bow_arrow,tag=hit,nbt={inGround:1b}] run data remove entity @s item.components."minecraft:potion_contents"
execute as @e[type=arrow,tag=maf_bow_arrow,tag=hit,nbt={inGround:1b}] run tag @s remove hit
execute as @e[type=arrow,tag=maf_bow_arrow,tag=ground,nbt={inGround:1b}] at @s run function maf:bow/run_ground with entity @s item.components."minecraft:custom_data".maf
