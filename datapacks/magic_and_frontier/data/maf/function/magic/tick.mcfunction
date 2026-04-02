# tell @a magic/tick

# クールダウンを1tick減算（下限0）
scoreboard players remove @a[scores={mafCoolTime=1..}] mafCoolTime 1

# キャスト中なら、キャスト処理を実行する
execute as @a at @s if score @s mafCastTime matches 0.. run function maf:magic/cast/tick

execute as @a at @s run function maf:magic/mp/mp_manage

function maf:magic/mp/mpbar
