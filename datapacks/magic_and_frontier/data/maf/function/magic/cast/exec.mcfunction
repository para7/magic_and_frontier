tellraw @s [{"text":"cast/exec"}]

# MPの消費処理
scoreboard players operation @s mafMP -= @s mafCastCost

# execの方に処理を回す
function #oh_my_dat:please
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting{kind:"grimoire"} run function maf:magic/cast/run_grimoire_effect with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting
execute if data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting{kind:"passive"} run function maf:magic/cast/run_passive_apply with storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting

# デバッグ用: storage の中身を tell で表示
# tellraw @s ["casting: ",{"nbt":"_[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting","storage":"oh_my_dat:","interpret":false}]

data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting

