tellraw @s [{"text":"詠唱が中断されました"}]
# TODO: ログアウト時に詠唱を中断させる
scoreboard players set @s mafCastCost 0
scoreboard players set @s mafCastTime -1
scoreboard players set @s mafCastTimeMax 0
function #oh_my_dat:please
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].magic.casting
