# MP自然回復処理
# execute as @a[scores={mafCastTime=..-1}] run scoreboard players add @s mafMPTick 1
# scoreboard players add @a[scores={mafCastTime=..-1}] mafMPTick 1

# 1秒に回復する内部値
scoreboard players add @s mafMPTick 10

# 回復・再初期化処理
execute if score @s mafMPTick matches 600.. run scoreboard players add @s mafMP 1
execute if score @s mafMPTick matches 600.. run scoreboard players set @s mafMPTick 0

# 装備由来の最大MPを計算
function maf:magic/mp/calc_equipment_maxmp

# 最大MP = min(ソウル, 装備MP)
scoreboard players operation @s mafMaxMP = @s mafSoul
scoreboard players operation @s mafMaxMP < @s mafEquipMP

# MPキャップ処理
scoreboard players operation @s mafMP < @s mafMaxMP
