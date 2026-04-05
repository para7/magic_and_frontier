# すべての変数を削除し、初期化する。
# 新規環境へのインストールを再現するデバッグ用

scoreboard objectives remove mafPlayerID
scoreboard objectives remove mafMoved
scoreboard objectives remove mafWalkCM
scoreboard objectives remove mafCoolTime
scoreboard objectives remove p7_setSkSlot
scoreboard objectives remove p7_setSkID
scoreboard objectives remove p7_setSkEnable
scoreboard objectives remove p7_skillSlot1
scoreboard objectives remove p7_skillSlot2
scoreboard objectives remove p7_skillSlot3
scoreboard objectives remove mafAviateCM
scoreboard objectives remove mafClimbCM
scoreboard objectives remove mafCrouchCM
scoreboard objectives remove mafFallCM
scoreboard objectives remove mafFlyCM
scoreboard objectives remove mafSprintCM
scoreboard objectives remove mafWalkWaterCM
scoreboard objectives remove mafUnderWaterCM
scoreboard objectives remove mafSwimCM
scoreboard objectives remove mafIsMovedY
# scoreboard objectives remove p7_posXpre
scoreboard objectives remove mafPosYpre
# scoreboard objectives remove p7_posZpre
# scoreboard objectives remove p7_posX
scoreboard objectives remove mafPosY
# scoreboard objectives remove p7_posZ
scoreboard objectives remove mafSoul
scoreboard objectives remove mafSoulTick
scoreboard objectives remove mafSoulReset
scoreboard objectives remove mafLogin
scoreboard objectives remove mafMP
scoreboard objectives remove mafMaxMP
scoreboard objectives remove mafCastCost
scoreboard objectives remove mafCastTime
scoreboard objectives remove mafCastTimeMax
scoreboard objectives remove mafMPTick
scoreboard objectives remove const0
scoreboard objectives remove tmp
scoreboard objectives remove tmp2

bossbar remove minecraft:mpbar1
bossbar remove minecraft:mpbar2
bossbar remove minecraft:mpbar3
bossbar remove minecraft:mpbar4
bossbar remove minecraft:mpbar5
bossbar remove minecraft:mpbar6
bossbar remove minecraft:mpbar7
bossbar remove minecraft:mpbar8
bossbar remove minecraft:mpbar9
bossbar remove minecraft:mpbar10
bossbar remove minecraft:mpbar11
bossbar remove minecraft:mpbar12
bossbar remove minecraft:mpbar13
bossbar remove minecraft:mpbar14
bossbar remove minecraft:mpbar15
bossbar remove minecraft:mpbar16
bossbar remove minecraft:mpbar17
bossbar remove minecraft:mpbar18
bossbar remove minecraft:mpbar19
bossbar remove minecraft:mpbar20

# load bootstrap is gated by this flag, so reinstall must clear it first.
data remove storage maf:runtime initialized

function maf:load
advancement revoke @a only maf:entered_world3
advancement revoke @a only maf:use_grimoire
# TODO: ここに仮置き 本来はログイン時毎回
function maf:system/set_player_id/run
