# MP自然回復処理
# execute as @a[scores={mafCastTime=..-1}] run scoreboard players add @s mafMPTick 1
# scoreboard players add @a[scores={mafCastTime=..-1}] mafMPTick 1

# 1秒に回復する内部値
scoreboard players add @a mafMPTick 10

# 回復・再初期化処理
execute as @a[scores={mafMPTick=600..}] run scoreboard players add @s mafMP 1
execute as @a[scores={mafMPTick=600..}] run scoreboard players set @s mafMPTick 0

# 最大MPをソウルと同値にする
execute as @a run scoreboard players operation @s mafMaxMP = @s mafSoul

# MPキャップ処理
execute as @a run scoreboard players operation @s mafMP < @s mafMaxMP


