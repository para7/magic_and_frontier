# tell @a magic/tick


# 杖の使用判定
execute as @a at @s if entity @s[scores={mafUseWand=1..,mafCastTime=..-1}] run function maf:magic/check

# 魔法の実行
execute as @a at @s if entity @s[scores={mafCastID=1..}] run function maf:magic/exec/tick

# キャスト中なら、キャスト処理を実行する
execute as @a at @s if score @s mafCastTime matches 0.. run function maf:magic/cast/tick

execute as @a at @s run function maf:magic/mp/mp_manage

function maf:magic/mp/mpbar

scoreboard players set @a mafCastID 0