# tickから呼ばれるので割り当て不要
# function #oh_my_dat:please

data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand

data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand.id set from entity @s SelectedItem.components."minecraft:custom_data".maf.passiveId
execute if data entity @s SelectedItem.components."minecraft:custom_data".maf.passiveSlot run data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand.slot set from entity @s SelectedItem.components."minecraft:custom_data".maf.passiveSlot
execute if data entity @s SelectedItem.components."minecraft:custom_data".maf.passiveCondition run data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand.condition set from entity @s SelectedItem.components."minecraft:custom_data".maf.passiveCondition

# always
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand{condition:"always"} run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand

# attack
execute if score @s mafMeleeHit matches 1.. if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand{condition:"attack"} run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand

# none: skip
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.mainhand
