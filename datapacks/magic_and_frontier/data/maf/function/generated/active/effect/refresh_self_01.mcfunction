effect clear @s minecraft:poison
effect clear @s minecraft:slowness
effect clear @s minecraft:mining_fatigue
effect clear @s minecraft:blindness
effect clear @s minecraft:wither
effect clear @s minecraft:hunger
effect clear @s minecraft:levitation
particle minecraft:happy_villager ~ ~1 ~ 0.3 0.3 0.3 1 6
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.8
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は リフレッシュ を唱えた！"}]
