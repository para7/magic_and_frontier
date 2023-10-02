# MP自然回復処理
# execute as @a[scores={p7_castTime=..-1}] run scoreboard players add @s p7_MPTick 1
# scoreboard players add @a[scores={p7_castTime=..-1}] p7_MPTick 1

# 1秒に回復する内部値
scoreboard players add @a p7_MPTick 10

# 回復・再初期化処理
execute as @a[scores={p7_MPTick=600..}] run scoreboard players add @s p7_MP 1
execute as @a[scores={p7_MPTick=600..}] run scoreboard players set @s p7_MPTick 0

# 最大MPをソウルと同値にする
execute as @a run scoreboard players operation @s p7_MaxMP = @s p7_soul

# MPキャップ処理
execute as @a run scoreboard players operation @s p7_MP < @s p7_MaxMP


