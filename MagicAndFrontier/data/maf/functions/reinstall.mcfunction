scoreboard objectives remove p7_useWand
scoreboard objectives remove p7_playerID
scoreboard objectives remove p7_move
scoreboard objectives remove p7_walkCM
scoreboard objectives remove p7_aviateCM
scoreboard objectives remove p7_climbCM
scoreboard objectives remove p7_crouchCM
scoreboard objectives remove p7_fallCM
scoreboard objectives remove p7_flyCM
scoreboard objectives remove p7_sprintCM
scoreboard objectives remove p7_walkWaterCM
scoreboard objectives remove p7_underWaterCM
scoreboard objectives remove p7_swimCM
scoreboard objectives remove p7_isMovedY
scoreboard objectives remove p7_posXpre
scoreboard objectives remove p7_posYpre
scoreboard objectives remove p7_posZpre
scoreboard objectives remove p7_posX
scoreboard objectives remove p7_posY
scoreboard objectives remove p7_posZ
scoreboard objectives remove const0
scoreboard objectives remove tmp
scoreboard objectives remove tmp2

# # 


function maf:load
# advancement revoke @a only maf:entered_world
# ここに仮置き 本来はログイン時毎回
function maf:system/set_player_id/run