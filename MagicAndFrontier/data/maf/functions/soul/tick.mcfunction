# 1秒に回復する内部値
scoreboard players add @a p7_soulTick 10

# 回復・再初期化処理
execute as @a[scores={p7_soulTick=1200..}] run scoreboard players add @s p7_soul 1
execute as @a[scores={p7_soulTick=1200..}] run scoreboard players set @s p7_soulTick 0

# ソウル最大値
# キャップ処理
execute as @a[scores={p7_soul=101..}] run scoreboard players set @s p7_soul 100

execute as @a[scores={p7_soulReset=1..}] run scoreboard players set @s p7_soul 0
# MP回復も一定時間ストップ
execute as @a[scores={p7_soulReset=1..}] run scoreboard players set @s p7_MPTick 0
execute as @a[scores={p7_soulReset=1..}] run scoreboard players set @s p7_soulReset 0
