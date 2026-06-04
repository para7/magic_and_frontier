effect give @s minecraft:instant_health 1 1
execute at @s run particle minecraft:heart ~ ~1.0 ~ 0.3 0.5 0.3 1 8 force
playsound minecraft:entity.player.levelup master @a[distance=..24] ~ ~ ~ 1.5 2.0
