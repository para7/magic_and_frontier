function #oh_my_dat:please

# パッシブスロット（condition 確認は run_effect に委譲）
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1.id run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2.id run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot3.id run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot3


# メインハンド装備のパッシブ
# tmp を slot と同じ構造で初期化（condition デフォルト: always）
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp.id set from entity @s SelectedItem.components."minecraft:custom_data".maf.passiveId
execute if data entity @s SelectedItem.components."minecraft:custom_data".maf.passiveCondition run data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp.condition set from entity @s SelectedItem.components."minecraft:custom_data".maf.passiveCondition

# condition 確認は run_effect に委譲
function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp

data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.tmp


# 弓着弾処理
execute if score @s mafBowHit matches 1.. run function maf:bow/on_bow_hit
