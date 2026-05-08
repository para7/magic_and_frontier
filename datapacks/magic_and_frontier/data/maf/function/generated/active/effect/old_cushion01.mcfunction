particle minecraft:smoke ~ ~ ~ 1.2 1 1.2 1 200
playsound minecraft:block.anvil.use player @s ~ ~ ~ 0.8 2.0
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5
execute if entity @s[nbt={Inventory:[{id:"minecraft:slime_block"}]}] run setblock ~ ~-1 ~ minecraft:slime_block keep
execute if entity @s[nbt={Inventory:[{id:"minecraft:slime_block"}]}] run clear @s minecraft:slime_block 1
tellraw @a[distance=..10] [{"selector":"@s"},{"text":" は クッション を唱えた！"}]
