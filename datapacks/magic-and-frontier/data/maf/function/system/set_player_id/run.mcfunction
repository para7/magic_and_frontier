# 通し番号を付与していきます 1スタート
data modify storage p7:maf playerSequence set value 0b
execute as @a run function maf:system/set_player_id/do_not_call