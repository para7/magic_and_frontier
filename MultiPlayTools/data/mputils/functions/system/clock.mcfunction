
# scoreboard players operation @s Dummy = @s p7PlayTime
# scoreboard players operation @s Dummy /= @s const20
# scoreboard players operation @s p7PTSeconds = @s Dummy

# scoreboard players operation @s Dummy = @s p7PTSeconds
# scoreboard players operation @s Dummy /= @s const60
# scoreboard players operation @s p7PTMinutes = @s Dummy

scoreboard players operation @s Dummy = @s p7PlayTime
scoreboard players operation @s Dummy /= @s const1200
scoreboard players operation @s p7PTMinutes = @s Dummy

scoreboard players operation @s Dummy = @s p7PTMinutes
scoreboard players operation @s Dummy /= @s const60
scoreboard players operation @s p7PTHours = @s Dummy