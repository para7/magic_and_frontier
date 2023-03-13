tellraw @s [{"text":"cast/exec id:"}, {"score":{"name":"@s","objective":"p7_castID"}}]

# MPの消費処理
scoreboard players operation @s p7_MP -= @s p7_castCost

execute if entity @s[scores={p7_castID=1}] run function maf:magic/effect/magic1