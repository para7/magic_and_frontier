# アクティブの仕様判定を行う処理
# アイテム指定はしてないのであらゆるアイテムで処理は呼び出される
advancement revoke @s only maf:use_active

execute unless entity @s[scores={mafCastTime=..-1}] run return 0
scoreboard players add @s mafCoolTime 0
execute unless entity @s[scores={mafCoolTime=..0}] run return 0

function #oh_my_dat:please
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting
data remove storage maf:tmp magic_spell_loader

execute if data entity @s SelectedItem.components."minecraft:custom_data".maf.active_id run data modify storage maf:tmp magic_spell_loader.active.id set from entity @s SelectedItem.components."minecraft:custom_data".maf.active_id
execute if data storage maf:tmp magic_spell_loader.active.id run function maf:magic/exec/load_active_spell with storage maf:tmp magic_spell_loader.active

execute unless data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting if data entity @s SelectedItem.components."minecraft:custom_data".maf.passive run data modify storage maf:tmp magic_spell_loader.passive set from entity @s SelectedItem.components."minecraft:custom_data".maf.passive
execute unless data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting if data storage maf:tmp magic_spell_loader.passive.id if data storage maf:tmp magic_spell_loader.passive.slot run function maf:magic/exec/load_passive_spell with storage maf:tmp magic_spell_loader.passive

execute unless data storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].maf.magic.casting run return 0
function maf:magic/exec/set_magic
