# buff_queue を先頭から1つ取り出して処理する再帰関数。
# tick.mcfunction が buff → buff_queue にコピーしてからこの関数を呼ぶ。
# 残存バフは buff_current を buff に append して次ティックへ引き継ぐ。

# キューがなければ終了
execute unless data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_queue[0] run return 0

# buff_queue の先頭を buff_current に移動
data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current set from storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_queue[0]
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_queue[0]

# tick_function が定義されていればバフ効果を実行
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current.tick_function run function maf:buff/run_tick with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current

# 残りティックを1減算して tmp に保存
execute store result score @s tmp run data get storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current.tick
scoreboard players remove @s tmp 1
execute store result storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current.tick int 1 run scoreboard players get @s tmp

# tick が 0 以下 → 期限切れ。destructor があれば実行して破棄
execute if score @s tmp matches ..0 if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current.destructor_function run function maf:buff/run_destructor with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current
# tick が 1 以上 → 次ティックへ継続
execute if score @s tmp matches 1.. run data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff append from storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current

data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.buff_current

# 次のキュー要素を処理（再帰）
function maf:buff/process_queue
