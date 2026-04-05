# グリモアの仕様判定を行う処理
# アイテム指定はしてないのであらゆるアイテムで処理は呼び出される
advancement revoke @s only maf:use_grimoire

execute unless entity @s[scores={mafCastTime=..-1}] run return 0
scoreboard players add @s mafCoolTime 0
execute unless entity @s[scores={mafCoolTime=..0}] run return 0

function #oh_my_dat:please
data remove storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].magic.casting
execute unless data entity @s SelectedItem.components."minecraft:custom_data".maf.spell run return 0

data modify storage oh_my_dat: _[-4][-4][-4][-4][-4][-4][-4][-4].magic.casting set from entity @s SelectedItem.components."minecraft:custom_data".maf.spell
function maf:magic/exec/set_magic
