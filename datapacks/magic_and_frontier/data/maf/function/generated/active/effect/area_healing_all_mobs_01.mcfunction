execute as @e[distance=..10] run effect give @s minecraft:instant_health 1 0
execute at @e[distance=..10] run particle minecraft:heart ~ ~1 ~ 0.3 0.3 0.3 1 4
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 2.0
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は エリアヒーリング を唱えた！"}]
