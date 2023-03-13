tellraw @s [{"text":"使用した魔法ID : "},{"score":{"name":"@s","objective":"p7_magicID"}}]

# 初期化
data remove storage p7:maf magictmp

# 該当データを一時変数にロードする
execute if entity @s[scores={p7_magicID=1}] run data modify storage p7:maf magictmp set from storage p7:maf_magicdb data.m1

execute if entity @s[scores={p7_magicID=2}] run data modify storage p7:maf magictmp set from storage p7:maf_magicdb data.m2

# データがあればキャスト処理に移る
execute if data storage p7:maf magictmp run function maf:magic/cast/set_magic 