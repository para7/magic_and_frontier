# MP自然回復処理
# execute as @a[scores={p7_castTime=..-1}] run scoreboard players add @s p7_MPTick 1
# scoreboard players add @a[scores={p7_castTime=..-1}] p7_MPTick 1
scoreboard players add @a p7_MPTick 1
execute as @a[scores={p7_MPTick=20..}] run scoreboard players add @s p7_MP 1
execute as @a[scores={p7_MPTick=20..}] run scoreboard players set @s p7_MPTick 0

# MPキャップ処理
execute as @a run scoreboard players operation @s p7_MP < @s p7_MaxMP