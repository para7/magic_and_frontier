function #oh_my_dat:please
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2.id set value "passive_1"
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2.condition set value "always"
tellraw @s [{"text":"[slot2]に[Quickstep]を設定しました"}]
