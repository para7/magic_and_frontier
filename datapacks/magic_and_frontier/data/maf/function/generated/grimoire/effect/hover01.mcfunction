
effect give @a[distance=..7] minecraft:levitation 17 255
# effect give @a[distance=..4] minecraft:speed 24 0


particle minecraft:cloud ~0 ~0.5 ~ 0.5 0.8 0.5 0.2 50 normal

playsound minecraft:entity.ender_dragon.ambient player @a ~ ~ ~ 0.8 1.5

playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5

tellraw @a[distance=..7] [{"selector":"@s"},{"text":" は ホバリング を唱えた！"}]