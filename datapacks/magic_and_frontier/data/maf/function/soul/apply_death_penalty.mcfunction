# リスポーン後に死亡ペナルティを確定する。
scoreboard players set @s tmp 10
scoreboard players operation @s mafSoul *= @s tmp
scoreboard players set @s tmp 100
scoreboard players operation @s mafSoul /= @s tmp

scoreboard players set @s mafMPTick 0

# 最大MPを即時再計算して、MP表示と現在MPを死亡後のソウル量に揃える。
function maf:magic/mp/calc_equipment_maxmp
scoreboard players operation @s mafMaxMP = @s mafEquipMP
scoreboard players operation @s mafMaxMP *= @s mafSoul
scoreboard players set @s tmp 100
scoreboard players operation @s mafMaxMP /= @s tmp
execute if score @s mafMaxMP matches 1000.. run scoreboard players set @s mafMaxMP 999
scoreboard players operation @s mafMP < @s mafMaxMP

scoreboard players set @s mafSoulReset 0
