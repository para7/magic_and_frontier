# スコアボードの前処理

execute as @a run function maf:prescore


# 処理開始

function maf:magic/tick

# 自動加算はここで消す
# dummyは各tick内で消す
scoreboard players set @a p7_useWand 0




scoreboard players set @a p7_walkCM 0
scoreboard players set @a p7_aviateCM 0
scoreboard players set @a p7_climbCM 0
scoreboard players set @a p7_crouchCM 0
scoreboard players set @a p7_fallCM 0
scoreboard players set @a p7_flyCM 0
scoreboard players set @a p7_sprintCM 0
scoreboard players set @a p7_swimCM 0
scoreboard players set @a p7_walkWaterCM 0
scoreboard players set @a p7_underWaterCM 0

function maf:test/tick
