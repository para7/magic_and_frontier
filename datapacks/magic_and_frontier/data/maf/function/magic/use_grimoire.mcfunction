# グリモアの仕様判定を行う処理
# アイテム指定はしてないのであらゆるアイテムで処理は呼び出される
advancement revoke @s only maf:use_grimoire

execute unless entity @s[scores={mafCastTime=..-1}] run return 0

data remove storage p7:maf magictmp
execute unless data entity @s SelectedItem.components."minecraft:custom_data".maf.spell run return 0

data modify storage p7:maf magictmp set from entity @s SelectedItem.components."minecraft:custom_data".maf.spell
function maf:magic/exec/set_magic
