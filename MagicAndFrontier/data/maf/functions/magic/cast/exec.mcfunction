tellraw @s [{"text":"cast/exec id:"}, {"score":{"name":"@s","objective":"p7_castID"}}]

# MPの消費処理
scoreboard players operation @s mafMP -= @s p7_castCost

# execの方に処理を回す
function maf:magic/exec/selectexec
