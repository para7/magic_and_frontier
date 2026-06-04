effect give @s minecraft:instant_health 1 0
particle minecraft:heart ~ ~1 ~ 0.3 0.3 0.3 1 4
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 2.0
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は キュア を唱えた！"}]
