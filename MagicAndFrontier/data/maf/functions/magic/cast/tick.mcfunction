# tellraw @s [{"text":"cast"},{"score":{"name":"@s","objective":"mafCastTime"}}]

# ブタに詠唱を止められないようにする
effect give @e[type=minecraft:pig, distance=..4] minecraft:slowness 10 127 true

scoreboard players set @a mafMPTick -40

#  particle minecraft:enchant ~ ~1.5 ~ 0.5 0 0.5 0.3 10 force
# 詠唱演出
execute if score @s mafCastTime matches 40.. run particle minecraft:enchant ~ ~2.3 ~ 0 0 0 3 2 force
execute if score @s mafCastTime matches 40 run particle minecraft:enchant ~ ~2.3 ~ 0 0 0 20 800 force

# 詠唱中の移動キャンセル 
# 滑りうちのしきい値設定
execute if entity @s[scores={mafCastTime=11..,mafMoved=1..}] run function maf:magic/cast/cancel

# 要消費MPチェック

execute if entity @s[scores={mafCastTime=0}] run function maf:magic/cast/exec

scoreboard players remove @s[scores={mafCastTime=0..}] mafCastTime 1
scoreboard players set @s p7_setSkEnable -1