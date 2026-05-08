fill ~-5 ~-3 ~-5 ~5 ~1 ~5 minecraft:frosted_ice replace minecraft:water[level=0]
execute as @e[distance=..12] at @s if block ~ ~-1 ~ minecraft:frosted_ice run damage @s 6 minecraft:freeze
execute as @e[distance=..12] at @s if block ~ ~ ~ minecraft:frosted_ice run damage @s 10 minecraft:freeze
execute as @e[distance=..12] at @s if block ~ ~-1 ~ minecraft:ice run damage @s 10 minecraft:freeze
execute as @e[distance=..12] at @s if block ~ ~-1 ~ minecraft:packed_ice run damage @s 16 minecraft:freeze
execute as @e[distance=..12] at @s if block ~ ~-1 ~ minecraft:blue_ice run damage @s 24 minecraft:freeze
tellraw @a[distance=..50] [{"selector":"@s"},{"text":" は アブソリュート・ゼロ を唱えた！"}]
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 2
playsound minecraft:entity.evoker.cast_spell master @a ~ ~ ~ 2 0.5
playsound minecraft:block.glass.break master @a ~ ~ ~ 3 0.6
