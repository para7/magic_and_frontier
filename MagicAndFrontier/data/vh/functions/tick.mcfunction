# function vh:replacer/tick

# 発生したばかりの敵に属性付与
execute as @e[type=#p7b:enemymob] unless entity @s[tag=enemymob] at @s run function vh:replacer/tick

# 修正済みタグを追加
tag @e[type=#p7b:enemymob] add vh_modified
