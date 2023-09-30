tellraw @a [{"text":"enable datapack: Magic and Frontier"}]

scoreboard objectives add p7_useWand minecraft.used:minecraft.carrot_on_a_stick

function maf:magic/load

scoreboard objectives add p7_playerID dummy

scoreboard objectives add p7_move dummy


scoreboard objectives add p7_walkCM minecraft.custom:minecraft.walk_one_cm

# スキルセット処理
# セット先スロット
scoreboard objectives add p7_setSkSlot trigger
# セットID
scoreboard objectives add p7_setSkID dummy
# triggerを無効化出来ないので、有効判定用
scoreboard objectives add p7_setSkEnable dummy

scoreboard objectives add p7_skillSlot1 dummy
scoreboard objectives add p7_skillSlot2 dummy
scoreboard objectives add p7_skillSlot3 dummy


# エリトラ
scoreboard objectives add p7_aviateCM minecraft.custom:minecraft.aviate_one_cm
# はしご・つた
scoreboard objectives add p7_climbCM minecraft.custom:minecraft.climb_one_cm
# スニーク
scoreboard objectives add p7_crouchCM minecraft.custom:minecraft.crouch_one_cm
# 落下距離 1ブロックから
scoreboard objectives add p7_fallCM minecraft.custom:minecraft.fall_one_cm
# 空中移動 1ブロックから
scoreboard objectives add p7_flyCM minecraft.custom:minecraft.fly_one_cm
# スプリント
scoreboard objectives add p7_sprintCM minecraft.custom:minecraft.sprint_one_cm
# 水面移動
scoreboard objectives add p7_walkWaterCM minecraft.custom:minecraft.walk_on_water_one_cm
# 水中立ち泳ぎ
scoreboard objectives add p7_underWaterCM minecraft.custom:minecraft.walk_under_water_one_cm
# 泳ぎ
scoreboard objectives add p7_swimCM minecraft.custom:minecraft.swim_one_cm

# Y座標に変化があったら1が入る
scoreboard objectives add p7_isMovedY dummy

# scoreboard objectives add p7_posXpre dummy
# scoreboard objectives add p7_posYpre dummy
# scoreboard objectives add p7_posZpre dummy

# scoreboard objectives add p7_posX dummy
scoreboard objectives add p7_posY dummy
# scoreboard objectives add p7_posZ dummy

# # ソウルシステム用
scoreboard objectives add p7_soul dummy "ソウル"
scoreboard objectives add p7_soulTick dummy
scoreboard objectives add p7_soulReset deathCount 


# 定数
# scoreboard objectives add const0 dummy

# 諸計算用
scoreboard objectives add tmp dummy
scoreboard objectives add tmp2 dummy

gamerule keepInventory true