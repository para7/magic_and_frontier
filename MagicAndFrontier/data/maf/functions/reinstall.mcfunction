scoreboard objectives remove mafUseWand
scoreboard objectives remove mafPlayerID
scoreboard objectives remove mafMoved
scoreboard objectives remove mafWalkCM
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
scoreboard objectives remove const0
scoreboard objectives remove tmp
scoreboard objectives remove tmp2

# # 


function maf:load
advancement revoke @a only maf:entered_world2
# ここに仮置き 本来はログイン時毎回
function maf:system/set_player_id/run