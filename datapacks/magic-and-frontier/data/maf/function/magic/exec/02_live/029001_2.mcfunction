
function maf:magic/exec/selectdb

tellraw @s [{"nbt":"magictmp.title","storage":"p7:maf"}]
tellraw @s [{"text":"詠唱tick: "},{"nbt":"magictmp.cast","storage":"p7:maf"}]
tellraw @s [{"text":"消費MP : "},{"nbt":"magictmp.cost","storage":"p7:maf"}]
tellraw @s [{"text":"効果: "},{"nbt":"magictmp.description","storage":"p7:maf"}]

# modifires 置き換え
function maf:magic/exec/generated/analyze

