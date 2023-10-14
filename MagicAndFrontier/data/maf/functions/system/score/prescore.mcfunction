# Y移動検知
scoreboard players operation @s tmp = @s mafPosY
scoreboard players set @s mafIsMovedY 0

# execute store result score @s p7_posX run data get entity @s Pos[0]
execute store result score @s mafPosY run data get entity @s Pos[1]
# execute store result score @s p7_posZ run data get entity @s Pos[2]

execute unless score @s tmp = @s mafPosY run scoreboard players set @s mafIsMovedY 1

scoreboard players set @s mafMoved 0
scoreboard players operation @s mafMoved += @s mafWalkCM
scoreboard players operation @s mafMoved += @s mafAviateCM
scoreboard players operation @s mafMoved += @s mafClimbCM
scoreboard players operation @s mafMoved += @s mafCrouchCM
scoreboard players operation @s mafMoved += @s mafFallCM
scoreboard players operation @s mafMoved += @s mafFlyCM
scoreboard players operation @s mafMoved += @s mafSprintCM
scoreboard players operation @s mafMoved += @s mafSwimCM
scoreboard players operation @s mafMoved += @s mafWalkWaterCM
scoreboard players operation @s mafMoved += @s mafUnderWaterCM
# scoreboard players operation @s mafMoved += @s mafIsMovedY