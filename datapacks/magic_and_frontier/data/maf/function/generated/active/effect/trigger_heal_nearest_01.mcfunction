execute as @e[type=#maf:friendmob,distance=0.1..10,sort=nearest,limit=1] run function maf:buff/data/trigger_heal/trigger_heal_01/init
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は トリガーヒール を唱えた！"}]
