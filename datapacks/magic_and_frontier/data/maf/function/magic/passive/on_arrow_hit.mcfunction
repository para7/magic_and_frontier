advancement revoke @s only maf:arrow_hit

# Tagged arrow must exist near the hit target.
execute unless entity @e[type=arrow,tag=maf_passive_arrow,limit=1,sort=nearest,distance=..10] run return 0


# Cache passive metadata before removing the arrow.
data modify storage maf:tmp arrow_passive set from entity @e[type=arrow,tag=maf_passive_arrow,sort=nearest,limit=1] item.components."minecraft:custom_data".maf

# Run bow passive script on the damaged entity marked by the temporary effect.
# execute as @e[nbt={active_effects:[{id:"minecraft:dolphins_grace",amplifier:80b}]}] at @s run function maf:magic/passive/run_bow_effect with storage maf:tmp arrow_passive
execute as @e[nbt={HurtTime:10s}] at @s run function maf:magic/passive/run_bow_effect with storage maf:tmp arrow_passive

# Remove the temporary marker effect.
effect clear @e[nbt={active_effects:[{id:"minecraft:dolphins_grace",amplifier:80b}]}] minecraft:dolphins_grace

# Reproduce the normal arrow cleanup after hit.
kill @e[type=arrow,tag=maf_passive_arrow,sort=nearest,limit=1]

# advancement revoke @s only maf:arrow_hit

data remove storage maf:tmp arrow_passive
