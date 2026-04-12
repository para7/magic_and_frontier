$execute as @e[type=arrow,tag=maf_bow_arrow,nbt={item:{components:{"minecraft:custom_data":{maf:{shooterPlayerID:$(bow_player_id)}}}}}] at @s run function maf:magic/bow/resolve_hit_arrow
