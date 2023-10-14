tellraw @a [{"text":"enable datapack: Magic and Frontier"}]

scoreboard objectives add mafUseWand minecraft.used:minecraft.carrot_on_a_stick

function maf:magic/load

scoreboard objectives add mafPlayerID dummy

scoreboard objectives add mafMoved dummy


scoreboard objectives add mafWalkCM minecraft.custom:minecraft.walk_one_cm

# # スキルセット処理
# # セット先スロット
# scoreboard objectives add p7_setSkSlot trigger
# # セットID
# scoreboard objectives add p7_setSkID dummy
# # triggerを無効化出来ないので、有効判定用
# scoreboard objectives add p7_setSkEnable dummy

# scoreboard objectives add p7_skillSlot1 dummy
# scoreboard objectives add p7_skillSlot2 dummy
# scoreboard objectives add p7_skillSlot3 dummy


# エリトラ
scoreboard objectives add mafAviateCM minecraft.custom:minecraft.aviate_one_cm
# はしご・つた
scoreboard objectives add mafClimbCM minecraft.custom:minecraft.climb_one_cm
# スニーク
scoreboard objectives add mafCrouchCM minecraft.custom:minecraft.crouch_one_cm
# 落下距離 1ブロックから
scoreboard objectives add mafFallCM minecraft.custom:minecraft.fall_one_cm
# 空中移動 1ブロックから
scoreboard objectives add mafFlyCM minecraft.custom:minecraft.fly_one_cm
# スプリント
scoreboard objectives add mafSprintCM minecraft.custom:minecraft.sprint_one_cm
# 水面移動
scoreboard objectives add mafWalkWaterCM minecraft.custom:minecraft.walk_on_water_one_cm
# 水中立ち泳ぎ
scoreboard objectives add mafUnderWaterCM minecraft.custom:minecraft.walk_under_water_one_cm
# 泳ぎ
scoreboard objectives add mafSwimCM minecraft.custom:minecraft.swim_one_cm

# Y座標に変化があったら1が入る
scoreboard objectives add mafIsMovedY dummy

# scoreboard objectives add p7_posXpre dummy
# scoreboard objectives add mafPosYpre dummy
# scoreboard objectives add p7_posZpre dummy

# scoreboard objectives add p7_posX dummy
scoreboard objectives add mafPosY dummy
# scoreboard objectives add p7_posZ dummy

# # ソウルシステム用
scoreboard objectives add mafSoul dummy "ソウル"
scoreboard objectives add mafSoulTick dummy
scoreboard objectives add mafSoulReset deathCount 


# 定数
# scoreboard objectives add const0 dummy

# 諸計算用
scoreboard objectives add tmp dummy
scoreboard objectives add tmp2 dummy

gamerule keepInventory true

# MPシステム初期化用
scoreboard objectives add mafLogin minecraft.custom:minecraft.leave_game
