# tellraw @s [{"text":"cast"},{"score":{"name":"@s","objective":"p7_castTime"}}]

# execute if score @s p7_const0 = @s p7_castTime run function maf:magic/cast/exec

scoreboard players set @a p7_MPTick -40

# 詠唱中の移動キャンセル 
execute if entity @s[scores={p7_castTime=9..,p7_move=1..}] run function maf:magic/cast/cancel

# 要消費MPチェック

execute if entity @s[scores={p7_castTime=0}] run function maf:magic/cast/exec

scoreboard players remove @s[scores={p7_castTime=0..}] p7_castTime 1
