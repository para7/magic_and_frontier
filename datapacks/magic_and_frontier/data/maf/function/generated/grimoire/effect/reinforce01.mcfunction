fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:smooth_stone replace minecraft:stone
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:smooth_stone_slab[type=bottom] replace minecraft:stone_slab[type=bottom]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:smooth_stone_slab[type=top] replace minecraft:stone_slab[type=top]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:smooth_stone_slab[type=double] replace minecraft:stone_slab[type=double]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:stone replace minecraft:cobblestone
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:stone_slab[type=bottom] replace minecraft:cobblestone_slab[type=bottom]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:stone_slab[type=top] replace minecraft:cobblestone_slab[type=top]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:stone_slab[type=double] replace minecraft:cobblestone_slab[type=double]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:blue_ice replace minecraft:packed_ice
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:packed_ice replace minecraft:ice
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:ice replace minecraft:frosted_ice
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:polished_andesite replace minecraft:andesite
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:polished_diorite replace minecraft:diorite
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:polished_granite replace minecraft:granite
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:polished_andesite_slab[type=bottom] replace minecraft:andesite_slab[type=bottom]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:polished_diorite_slab[type=bottom] replace minecraft:diorite_slab[type=bottom]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:polished_granite_slab[type=bottom] replace minecraft:granite_slab[type=bottom]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:polished_andesite_slab[type=top] replace minecraft:andesite_slab[type=top]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:polished_diorite_slab[type=top] replace minecraft:diorite_slab[type=top]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:polished_granite_slab[type=top] replace minecraft:granite_slab[type=top]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:polished_andesite_slab[type=double] replace minecraft:andesite_slab[type=double]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:polished_diorite_slab[type=double] replace minecraft:diorite_slab[type=double]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:polished_granite_slab[type=double] replace minecraft:granite_slab[type=double]
fill ~-14 ~-14 ~-14 ~14 ~14 ~14 minecraft:white_concrete replace minecraft:snow_block
particle minecraft:happy_villager ~ ~ ~ 10 10 10 1 3000
# fill ~ ~ ~ ~ ~ ~ minecraft:snow_block
playsound minecraft:entity.generic.explode player @a ~ ~ ~ 2 1
playsound minecraft:block.anvil.use player @s ~ ~ ~ 2 2.0
playsound minecraft:entity.evoker.cast_spell player @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell player @a ~ ~ ~ 2 0.5
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は リインフォース を唱えた！"}]
# xp add @s -1 levels
