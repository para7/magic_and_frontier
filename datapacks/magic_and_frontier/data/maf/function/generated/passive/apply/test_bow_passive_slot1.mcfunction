function #oh_my_dat:please
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1.id set value "test_bow_passive"
tellraw @s [{"text":"[slot1]に[Test Bow Passive]を設定しました"}]
