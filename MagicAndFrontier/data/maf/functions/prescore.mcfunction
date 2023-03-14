# Y移動検知
scoreboard players operation @s tmp = @s p7_posY
scoreboard players set @s p7_isMovedY 0

# execute store result score @s p7_posX run data get entity @s Pos[0]
execute store result score @s p7_posY run data get entity @s Pos[1]
# execute store result score @s p7_posZ run data get entity @s Pos[2]

execute unless score @s tmp = @s p7_posY run scoreboard players set @s p7_isMovedY 1

scoreboard players set @s p7_move 0
scoreboard players operation @s p7_move += @s p7_walkCM
scoreboard players operation @s p7_move += @s p7_aviateCM
scoreboard players operation @s p7_move += @s p7_climbCM
scoreboard players operation @s p7_move += @s p7_crouchCM
scoreboard players operation @s p7_move += @s p7_fallCM
scoreboard players operation @s p7_move += @s p7_flyCM
scoreboard players operation @s p7_move += @s p7_sprintCM
scoreboard players operation @s p7_move += @s p7_swimCM
scoreboard players operation @s p7_move += @s p7_walkWaterCM
scoreboard players operation @s p7_move += @s p7_underWaterCM
# scoreboard players operation @s p7_move += @s p7_isMovedY