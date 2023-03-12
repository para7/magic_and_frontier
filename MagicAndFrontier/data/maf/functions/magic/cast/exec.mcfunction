tellraw @s [{"text":"cast/exec id:"}, {"score":{"name":"@s","objective":"p7_castID"}}]

execute if entity @s[scores={p7_castID=1}] run function maf:magic/effect/magic1