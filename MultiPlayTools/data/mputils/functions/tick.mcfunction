execute if data storage para7sys:mpu {glow: true} run effect give @a minecraft:glowing 1 16 true

# プレイ時間計算
execute as @a run function mputils:system/clock

# ログインメッセージ
execute as @a[scores={p7logout=1..}] run function mputils:system/logintell

# ベッド通知
execute if data storage para7sys:mpu {bed: true} run function mputils:system/bednotice/main

# scoreboard objectives setdisplay sidebar p7PTSeconds
# scoreboard objectives setdisplay sidebar p7PTMinutes
# scoreboard objectives setdisplay sidebar p7PTHours