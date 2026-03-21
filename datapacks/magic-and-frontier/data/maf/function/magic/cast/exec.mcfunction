tellraw @s [{"text":"cast/exec id:"}, {"score":{"name":"@s","objective":"mafEffectID"}}]

# MPの消費処理
scoreboard players operation @s mafMP -= @s mafCastCost

# execの方に処理を回す
function maf:magic/exec/selectexec
