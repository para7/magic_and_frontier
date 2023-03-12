
# scoreboard players operation @s tmp = @s p7_PlayTime
# scoreboard players operation @s tmp /= @s const20
# scoreboard players operation @s p7_PTSeconds = @s tmp

# scoreboard players operation @s tmp = @s p7_PTSeconds
# scoreboard players operation @s tmp /= @s const60
# scoreboard players operation @s p7_PTMinutes = @s tmp

scoreboard players operation @s tmp = @s p7_PlayTime
scoreboard players operation @s tmp /= @s const1200
scoreboard players operation @s p7_PTMinutes = @s tmp

scoreboard players operation @s tmp = @s p7_PTMinutes
scoreboard players operation @s tmp /= @s const60
scoreboard players operation @s p7_PTHours = @s tmp