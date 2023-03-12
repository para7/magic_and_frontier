tellraw @a [{"text":"enable datapack: Magic and Frontier"}]

scoreboard objectives add p7_useWand minecraft.used:minecraft.carrot_on_a_stick

function maf:magic/load

# 諸計算用
scoreboard objectives add tmp dummy
scoreboard objectives add tmp2 dummy