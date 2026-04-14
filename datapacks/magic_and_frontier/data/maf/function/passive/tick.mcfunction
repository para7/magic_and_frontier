function #oh_my_dat:please

# パッシブスロット（condition 確認は run_effect に委譲）
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1.id run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2.id run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot3.id run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot3

# メインハンド装備のパッシブ
execute if data entity @s SelectedItem.components."minecraft:custom_data".maf.passiveId run function maf:passive/run_mainhand_effect

# 弓着弾処理
execute if score @s mafBowHit matches 1.. run function maf:bow/on_bow_hit
