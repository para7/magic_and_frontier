function #oh_my_dat:please
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1.id set value "passive_1"
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1.condition set value "always"
tellraw @s [{"text":"[slot1]に[passive_1]を設定しました"}]
