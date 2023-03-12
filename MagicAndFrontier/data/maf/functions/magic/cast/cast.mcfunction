# tellraw @s [{"text":"cast"}]

execute if score @s p7_castCost = @s p7_castTime run function maf:magic/cast/exec

scoreboard players remove @s[scores={p7_castTime=0..}] p7_castTime 1
