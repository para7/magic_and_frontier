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

# 最大MP = 装備MP * ソウル割合
scoreboard players operation @s mafMaxMP = @s mafEquipMP
scoreboard players operation @s mafMaxMP *= @s mafSoul
scoreboard players set @s tmp 1000
scoreboard players operation @s mafMaxMP /= @s tmp
execute if score @s mafMaxMP matches 1000.. run scoreboard players set @s mafMaxMP 999

# MPキャップ処理
scoreboard players operation @s mafMP < @s mafMaxMP
