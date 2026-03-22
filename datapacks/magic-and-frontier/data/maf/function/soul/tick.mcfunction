# 1秒に回復する内部値
scoreboard players add @a mafSoulTick 10

# 回復・再初期化処理
execute as @a[scores={mafSoulTick=1200..}] run scoreboard players add @s mafSoul 1
execute as @a[scores={mafSoulTick=1200..}] run scoreboard players set @s mafSoulTick 0

# ソウル最大値
# キャップ処理
execute as @a[scores={mafSoul=101..}] run scoreboard players set @s mafSoul 100

execute as @a[scores={mafSoulReset=1..}] run scoreboard players set @s mafSoul 0
# MP回復も一定時間ストップ
execute as @a[scores={mafSoulReset=1..}] run scoreboard players set @s mafMPTick 0
execute as @a[scores={mafSoulReset=1..}] run scoreboard players set @s mafSoulReset 0
