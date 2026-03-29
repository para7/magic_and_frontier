# tell @a magic/tick

# キャスト中なら、キャスト処理を実行する
execute as @a at @s if score @s mafCastTime matches 0.. run function maf:magic/cast/tick

execute as @a at @s run function maf:magic/mp/mp_manage

function maf:magic/mp/mpbar
