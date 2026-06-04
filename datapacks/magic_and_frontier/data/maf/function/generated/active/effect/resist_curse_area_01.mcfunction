execute as @e[type=#maf:friendmob,distance=..10] run function maf:buff/data/resist/resist_curse_01/init
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は レジストカース を唱えた！"}]
