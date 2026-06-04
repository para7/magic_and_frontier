function maf:buff/data/resist/resist_blind_01/init
playsound minecraft:entity.player.levelup player @a ~ ~ ~ 2 1.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は レジストブライン を唱えた！"}]
