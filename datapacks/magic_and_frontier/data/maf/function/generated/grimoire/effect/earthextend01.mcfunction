fill ~-5 ~-1 ~-5 ~5 ~-6 ~5 minecraft:stone replace minecraft:water
fill ~-5 ~-1 ~-5 ~5 ~-6 ~5 minecraft:stone replace minecraft:lava


particle minecraft:smoke ~ ~ ~ 1.2 0.1 1.2 1 200

playsound minecraft:block.anvil.use player @s ~ ~ ~ 0.8 2.0

playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5

tellraw @a[distance=..20] [{"selector":"@s"},{"text":" は アースエクステンド を唱えた！"}]
