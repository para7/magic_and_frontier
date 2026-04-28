effect give @e[type=#maf:enemymob,distance=..8] minecraft:poison 10 2
playsound minecraft:entity.witch.throw master @a ~ ~ ~ 1.5 0.8
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は ポイズン を唱えた！"}]
