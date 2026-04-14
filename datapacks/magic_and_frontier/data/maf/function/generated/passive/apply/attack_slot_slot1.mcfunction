function #oh_my_dat:please
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1.id set value "attack_slot"
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1.condition set value "attack"
tellraw @s [{"text":"[slot1]に[攻撃テストスロット]を設定しました"}]
