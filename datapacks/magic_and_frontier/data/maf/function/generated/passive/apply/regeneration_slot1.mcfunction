function #oh_my_dat:please
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].passive.slot1.id set value "regeneration"
tellraw @s [{"text":"[slot1]に[いつでもリジェネ]を設定しました"}]
