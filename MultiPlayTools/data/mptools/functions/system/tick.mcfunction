
# scoreboard players operation @s p7_mpuTmp = @s p7_PlayTime
# scoreboard players operation @s p7_mpuTmp /= @s const20
# scoreboard players operation @s p7_PTSeconds = @s p7_mpuTmp

# scoreboard players operation @s p7_mpuTmp = @s p7_PTSeconds
# scoreboard players operation @s p7_mpuTmp /= @s const60
# scoreboard players operation @s p7_PTMinutes = @s p7_mpuTmp

scoreboard players operation @s p7_mpuTmp = @s p7_PlayTime
scoreboard players operation @s p7_mpuTmp /= @s const1200
scoreboard players operation @s p7_PTMinutes = @s p7_mpuTmp

scoreboard players operation @s p7_mpuTmp = @s p7_PTMinutes
scoreboard players operation @s p7_mpuTmp /= @s const60
scoreboard players operation @s p7_PTHours = @s p7_mpuTmp