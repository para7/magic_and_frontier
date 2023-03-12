execute if data storage p7_:mpu {glow: true} run effect give @a minecraft:glowing 1 16 true

# プレイ時間計算
execute as @a run function mptools:system/tick

# ログインメッセージ
# execute as @a[scores={p7_logout=1..}] run function mptools:system/logintell

# ベッド通知
execute if data storage p7_:mpu {bed: true} run function mptools:system/bednotice/main

# scoreboard objectives setdisplay sidebar p7_PTSeconds
# scoreboard objectives setdisplay sidebar p7_PTMinutes
# scoreboard objectives setdisplay sidebar p7_PTHours