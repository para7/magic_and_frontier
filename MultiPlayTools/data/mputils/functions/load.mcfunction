# プレイヤー発光

# 定数
scoreboard objectives add const20 dummy
scoreboard objectives add const60 dummy
scoreboard objectives add const1200 dummy

# 画面表示系
scoreboard objectives add p7Health health "HP"
scoreboard objectives add p7DeathCnt deathCount "リスポーン数"
scoreboard objectives add p7PlayTime minecraft.custom:minecraft.play_time
# scoreboard objectives add p7PTSeconds dummy "プレイ時間(s)(未使用)"
scoreboard objectives add p7PTMinutes dummy "プレイ時間(m)"
scoreboard objectives add p7PTHours dummy "プレイ時間(h)"
scoreboard objectives add Dummy dummy

scoreboard objectives add p7logout minecraft.custom:minecraft.leave_game


scoreboard players set @s const1200 1200

scoreboard players set @s p7logout 0

# ベッド通知
scoreboard objectives add BedNotification minecraft.custom:minecraft.sleep_in_bed
data merge storage para7sys:mpu {bed: true}


tellraw @a [{"text":"enable datapack: MultiPlayTools"}]