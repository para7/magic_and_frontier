effect clear @s minecraft:blindness
particle minecraft:happy_villager ~ ~1 ~ 0.3 0.3 0.3 1 6
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.8
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は クリアブライン を唱えた！"}]
