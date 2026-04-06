function #oh_my_dat:please
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot1
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot2
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.passive.slot3
tellraw @s [{"text":"パッシブスロットをすべてクリアしました"}]
