tellraw @s [{"text":"cast/exec id:"}, {"score":{"name":"@s","objective":"p7_castID"}}]

# MPの消費処理
scoreboard players operation @s p7_MP -= @s p7_castCost

# 回復系
execute if entity @s[scores={p7_castID=1}] run function maf:magic/exec/00_heal/0001

# 攻撃系
execute if entity @s[scores={p7_castID=1001}] run function maf:magic/exec/01_attack/1001

# 生活系
execute if entity @s[scores={p7_castID=2001}] run function maf:magic/exec/02_live/2001

# デバフ系
execute if entity @s[scores={p7_castID=3001}] run function maf:magic/exec/03_debuff/3001

# 戦闘補助系
execute if entity @s[scores={p7_castID=4001}] run function maf:magic/exec/04_buff/4001

# スキル系
execute if entity @s[scores={p7_castID=10001}] run function maf:magic/exec/10_skill/selector