
# fill ~-2 ~-1 ~-2 ~-2 ~2 ~2 minecraft:cobblestone replace minecraft:air

# fill ~2 ~-1 ~-2 ~2 ~2 ~2 minecraft:cobblestone replace minecraft:air

# fill ~-2 ~-1 ~-2 ~2 ~2 ~-2 minecraft:cobblestone replace minecraft:air

# fill ~-2 ~-1 ~2 ~2 ~2 ~2 minecraft:cobblestone replace minecraft:air

fill ~-2 ~3 ~-2 ~2 ~3 ~2 minecraft:cobblestone replace #p7b:air

particle minecraft:smoke ~ ~3 ~ 1.2 1 1.2 1 200


playsound minecraft:block.anvil.use player @s ~ ~3 ~ 0.8 2.0

playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5

tellraw @a[distance=..20] [{"selector":"@s"},{"text":" は ルーフ を唱えた！"}]
