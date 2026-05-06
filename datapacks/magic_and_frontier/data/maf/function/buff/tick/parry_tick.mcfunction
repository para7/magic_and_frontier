execute at @s run particle minecraft:cloud ~ ~1.0 ~ 0.15 0.25 0.15 0.01 3 force
# 10s だと判定が取れなかった。
execute unless data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current{triggered:1b} if entity @s[nbt={HurtTime:9s}] at @s run playsound minecraft:block.anvil.land master @a[distance=..24] ~ ~ ~ 1.0 1.7
