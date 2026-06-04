execute as @e[type=#maf:friendmob,distance=0.1..10] run function maf:buff/data/delay_heal/delay_heal_01/init
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は ディレイヒール を唱えた！"}]
