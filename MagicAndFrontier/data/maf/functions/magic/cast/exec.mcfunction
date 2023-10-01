tellraw @s [{"text":"cast/exec id:"}, {"score":{"name":"@s","objective":"p7_castID"}}]

# MPの消費処理
scoreboard players operation @s p7_MP -= @s p7_castCost

function maf:magic/cast/selectexec
