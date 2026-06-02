scoreboard players set #maf_item_check_ver tmp2 0
$execute store result score #maf_item_check_ver tmp2 run data get entity @s $(equip).components."minecraft:custom_data".maf.ver
execute if score #maf_item_check_ver tmp2 = #maf_item_ver tmp run return 0

$data modify storage maf:tmp item_update.source_id set from entity @s $(equip).components."minecraft:custom_data".maf.source_id
$data modify storage maf:tmp item_update.slot set value "$(slot)"
$data modify storage maf:tmp item_update.equip set value "$(equip)"
function maf:item/update/dispatch with storage maf:tmp item_update
