# @s = arrow entity, macro args: passive_id, life
$data merge entity @s {Tags:["maf_passive_arrow"],PierceLevel:2b,item:{components:{"minecraft:custom_data":{maf:{passiveId:"$(passive_id)"}},"minecraft:potion_contents":{custom_effects:[{id:"minecraft:dolphins_grace",duration:200,amplifier:80}]}}}}
$data modify entity @s life set value $(life)s
data modify entity @s item.components."minecraft:custom_data".maf.shooterPlayerID set from storage maf:tmp bow_player_id
