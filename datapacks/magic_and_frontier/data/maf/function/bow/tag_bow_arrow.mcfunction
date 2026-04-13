# @s = arrow entity, args: bow_id, life
$data merge entity @s {Tags:["maf_bow_arrow","maf_bow_arrow_new"],item:{components:{"minecraft:custom_data":{maf:{bowId:"$(bow_id)"}}}}}
$data modify entity @s life set value $(life)s
data modify entity @s item.components."minecraft:custom_data".maf.shooterPlayerID set from storage maf:tmp bow_player_id
