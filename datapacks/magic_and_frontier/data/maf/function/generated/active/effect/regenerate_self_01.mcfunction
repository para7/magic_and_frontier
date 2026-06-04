effect give @s minecraft:regeneration 30 1
particle minecraft:heart ~ ~1 ~ 0.3 0.3 0.3 0.5 6
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 2.0
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は リジェネレート を唱えた！"}]
