fill ~-6 ~-1 ~-6 ~6 ~-1 ~6 minecraft:cobblestone replace #p7b:air

particle minecraft:smoke ~ ~ ~ 1.2 0.1 1.2 1 200

playsound minecraft:block.anvil.use player @s ~ ~ ~ 0.8 2.0
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5

tellraw @a[distance=..20] [{"selector":"@s"},{"text":" は "},{"nbt":"magictmp.title","storage":"p7:maf"},{"text": "を唱えた！"}]