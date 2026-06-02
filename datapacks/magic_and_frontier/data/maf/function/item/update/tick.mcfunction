function maf:generated/item/update/timestamp

execute if data entity @s SelectedItem.components."minecraft:custom_data".maf.source_id run function maf:item/update/check_slot {slot:"weapon.mainhand",equip:"SelectedItem"}
execute if data entity @s equipment.offhand.components."minecraft:custom_data".maf.source_id run function maf:item/update/check_slot {slot:"weapon.offhand",equip:"equipment.offhand"}
execute if data entity @s equipment.head.components."minecraft:custom_data".maf.source_id run function maf:item/update/check_slot {slot:"armor.head",equip:"equipment.head"}
execute if data entity @s equipment.chest.components."minecraft:custom_data".maf.source_id run function maf:item/update/check_slot {slot:"armor.chest",equip:"equipment.chest"}
execute if data entity @s equipment.legs.components."minecraft:custom_data".maf.source_id run function maf:item/update/check_slot {slot:"armor.legs",equip:"equipment.legs"}
execute if data entity @s equipment.feet.components."minecraft:custom_data".maf.source_id run function maf:item/update/check_slot {slot:"armor.feet",equip:"equipment.feet"}
