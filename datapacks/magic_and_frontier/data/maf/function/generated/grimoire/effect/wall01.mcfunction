fill ~-2 ~-1 ~-2 ~-2 ~2 ~2 minecraft:cobblestone replace #maf:air
fill ~2 ~-1 ~-2 ~2 ~2 ~2 minecraft:cobblestone replace #maf:air
fill ~-2 ~-1 ~-2 ~2 ~2 ~-2 minecraft:cobblestone replace #maf:air
fill ~-2 ~-1 ~2 ~2 ~2 ~2 minecraft:cobblestone replace #maf:air
particle minecraft:smoke ~ ~ ~ 1.2 1 1.2 1 200
playsound minecraft:block.anvil.use player @s ~ ~ ~ 0.8 2.0
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5
tellraw @a[distance=..20] [{"selector":"@s"},{"text":" は ウォール を唱えた！"}]
