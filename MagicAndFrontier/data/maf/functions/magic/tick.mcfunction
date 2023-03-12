# 杖の使用判定
execute if entity @s[scores={p7_useWand=1..}] run function maf:magic/check

execute if entity @s[scores={p7_magicID=1..}] run function maf:magic/exec

scoreboard players set @a p7_magicID 0