# tickから呼ばれるので割り当て不要
# function #oh_my_dat:please

# tmp を slot と同じ構造で初期化（condition デフォルト: always）
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp.id set from entity @s SelectedItem.components."minecraft:custom_data".maf.passiveId
execute if data entity @s SelectedItem.components."minecraft:custom_data".maf.passiveCondition run data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp.condition set from entity @s SelectedItem.components."minecraft:custom_data".maf.passiveCondition

# condition 確認は run_effect に委譲
function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp

data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp
