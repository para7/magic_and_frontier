# 杖の使用判定
execute as @a at @s if entity @s[scores={p7_useWand=1..,p7_castTime=..-1}] run function maf:magic/check

# 魔法の実行
execute as @a at @s if entity @s[scores={p7_magicID=1..}] run function maf:magic/exec/tick

# キャスト中なら、キャスト処理を実行する
execute as @a at @s if score @s p7_castTime matches 0.. run function maf:magic/cast/tick

execute as @a at @s run function maf:magic/mp_manage

function maf:magic/mpbar

scoreboard players set @a p7_magicID 0