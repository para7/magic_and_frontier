# 1秒に回復する内部値
execute as @a unless score @s mafSoulReset matches 1.. run scoreboard players add @s mafSoulTick 10

# 回復・再初期化処理
execute as @a[scores={mafSoulTick=1200..}] run scoreboard players add @s mafSoul 1
execute as @a[scores={mafSoulTick=1200..}] run scoreboard players set @s mafSoulTick 0

# ソウル最大値
# キャップ処理
execute as @a[scores={mafSoul=101..}] run scoreboard players set @s mafSoul 100

# 死亡画面で放置している間はMP/ソウル回復と詠唱を停止し、
# リスポーンしてHealthが戻ったtickで死亡ペナルティを適用する。
execute as @a[scores={mafSoulReset=1..}] run scoreboard players set @s mafMPTick 0
execute as @a[scores={mafSoulReset=1..}] run scoreboard players set @s mafCastCost 0
execute as @a[scores={mafSoulReset=1..}] run scoreboard players set @s mafCastTime -1
execute as @a[scores={mafSoulReset=1..}] run scoreboard players set @s mafCastTimeMax 0
execute as @a[scores={mafSoulReset=1..}] run scoreboard players set @s tmp 0
execute as @a[scores={mafSoulReset=1..}] store result score @s tmp run data get entity @s Health 1
execute as @a[scores={mafSoulReset=1..}] if score @s tmp matches 1.. run function maf:soul/apply_death_penalty
