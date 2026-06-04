execute as @e[type=#maf:friendmob,distance=0.1..10,sort=nearest,limit=1] run effect clear @s minecraft:mining_fatigue
execute at @e[type=#maf:friendmob,distance=0.1..10,sort=nearest,limit=1] run particle minecraft:happy_villager ~ ~1 ~ 0.3 0.3 0.3 1 6
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.8
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は クリアスロウ を唱えた！"}]
