function #oh_my_dat:please

# スキルスロット1〜3
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1.id run function maf:magic/passive/dispatch/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2.id run function maf:magic/passive/dispatch/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot3.id run function maf:magic/passive/dispatch/run_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot3

# メインハンド
execute if data entity @s SelectedItem.components."minecraft:custom_data".maf.spell{kind:"passive"} run function maf:magic/passive/dispatch/run_mainhand_effect
