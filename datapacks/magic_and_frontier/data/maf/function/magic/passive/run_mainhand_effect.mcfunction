data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand.id set from entity @s SelectedItem.components."minecraft:custom_data".maf.passiveId
execute if data entity @s SelectedItem.components."minecraft:custom_data".maf.passiveSlot run data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand.slot set from entity @s SelectedItem.components."minecraft:custom_data".maf.passiveSlot
execute if data entity @s SelectedItem.components."minecraft:custom_data".maf.passiveCondition run data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand.condition set from entity @s SelectedItem.components."minecraft:custom_data".maf.passiveCondition
function maf:magic/passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand
