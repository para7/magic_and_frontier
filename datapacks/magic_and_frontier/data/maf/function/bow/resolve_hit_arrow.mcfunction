execute unless entity @e[nbt={active_effects:[{id:"minecraft:dolphins_grace",amplifier:80b}]},sort=nearest,limit=1,distance=..5] run return 0
data modify storage maf:tmp bow_passive set from entity @s item.components."minecraft:custom_data".maf
execute if entity @s[tag=hit] as @e[nbt={active_effects:[{id:"minecraft:dolphins_grace",amplifier:80b}]},sort=nearest,limit=1,distance=..5] at @s run function maf:bow/run_bow_effect with storage maf:tmp bow_passive
effect clear @e[nbt={active_effects:[{id:"minecraft:dolphins_grace",amplifier:80b}]},sort=nearest,limit=1,distance=..5] minecraft:dolphins_grace
kill @s
data remove storage maf:tmp bow_passive
