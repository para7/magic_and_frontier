function #oh_my_dat:please

# always: 毎 tick 実行
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1{condition:"always"} run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2{condition:"always"} run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot3{condition:"always"} run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot3

# attack: 近接攻撃時のみ実行
execute if score @s mafMeleeHit matches 1.. if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1{condition:"attack"} run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1
execute if score @s mafMeleeHit matches 1.. if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2{condition:"attack"} run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2
execute if score @s mafMeleeHit matches 1.. if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot3{condition:"attack"} run function maf:passive/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot3

# none: 自動発動なし（他システムから function maf:generated/passive/effect/{id} で直接呼び出す）

# メインハンド装備のパッシブ
execute if data entity @s SelectedItem.components."minecraft:custom_data".maf.passiveId run function maf:passive/run_mainhand_effect

# 弓着弾処理
execute if score @s mafBowHit matches 1.. run function maf:bow/on_bow_hit
